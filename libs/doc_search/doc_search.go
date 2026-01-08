package doc_search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	TikaURL     string
	ESHost      string
	IndexName   string
	CallbackURL string

	MaxWorkers int
	QueueLimit int

	DefaultTopFragments int
	DefaultMaxContent   int
	FragmentSize        int
	DefaultMSM          string

	MaxPageSize     int
	MaxFrom         int
	DefaultPageSize int

	Port int
}

func GetEnv(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}

func getenvInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func LoadConfig() Config {
	return Config{
		TikaURL:     GetEnv("TIKA_URL", "http://192.168.12.49:9998/tika"),
		ESHost:      GetEnv("ES_HOST", "http://192.168.12.49:9200"),
		IndexName:   GetEnv("INDEX_NAME", "test_hikb2"),
		CallbackURL: GetEnv("CALLBACK_URL", "http://192.168.12.49:8001"),

		MaxWorkers: getenvInt("MAX_WORKERS", 4),
		QueueLimit: getenvInt("QUEUE_LIMIT", 20),

		DefaultTopFragments: getenvInt("DEFAULT_TOP_FRAGMENTS", 10),
		DefaultMaxContent:   getenvInt("DEFAULT_MAX_CONTENT", 80),
		FragmentSize:        getenvInt("FRAGMENT_SIZE", 300),
		DefaultMSM:          GetEnv("DEFAULT_MSM", "1"),

		MaxPageSize:     getenvInt("MAX_PAGE_SIZE", 100),
		MaxFrom:         getenvInt("MAX_FROM", 9900),
		DefaultPageSize: getenvInt("DEFAULT_PAGE_SIZE", 20),

		Port: getenvInt("FLASK_PORT", 10821),
	}
}

const (
	TagStart = `<span class="highlight term0">`
	TagEnd   = `</span>`
)

var wsRe = regexp.MustCompile(`\s+`)

func sanitizeText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return wsRe.ReplaceAllString(s, " ")
}

// 合并 </span>  ...  <span ...> 之间仅有空格或点号的情况
func mergeHighlightTags(s string) string {
	if s == "" {
		return s
	}
	pat := regexp.MustCompile(regexp.QuoteMeta(TagEnd) + `([\s\.]*)` + regexp.QuoteMeta(TagStart))
	return pat.ReplaceAllString(s, `$1`)
}

func splitKeepDelims(s string, delims []string) []string {
	type hit struct {
		pos int
		d   string
	}
	var res []string
	for len(s) > 0 {
		best := hit{pos: -1}
		for _, d := range delims {
			i := strings.Index(s, d)
			if i >= 0 && (best.pos == -1 || i < best.pos) {
				best = hit{pos: i, d: d}
			}
		}
		if best.pos == -1 {
			res = append(res, s)
			break
		}
		if best.pos > 0 {
			res = append(res, s[:best.pos])
		}
		res = append(res, best.d)
		s = s[best.pos+len(best.d):]
	}
	return res
}

// 与 Python 版本同思路：让高亮尽量靠前，长度受控且不破坏标签
func processHighlightContent(rawFragments []string, maxContent int) (string, bool) {
	var kept []string
	for _, f := range rawFragments {
		ff := strings.TrimSpace(f)
		if ff == "" {
			continue
		}
		if strings.Contains(ff, TagStart) {
			kept = append(kept, ff)
		}
	}
	if len(kept) == 0 {
		return "", false
	}
	combined := mergeHighlightTags(strings.Join(kept, "..."))
	if !strings.Contains(combined, TagStart) {
		return "", false
	}

	firstIdx := strings.Index(combined, TagStart)
	threshold := int(float64(maxContent) * 0.6)
	prefixBuf := int(float64(maxContent) * 0.4)

	if firstIdx > threshold {
		startPos := firstIdx - prefixBuf
		if startPos < 0 {
			startPos = 0
		}
		if startPos < len(combined) {
			combined = combined[startPos:]
		}
	}

	parts := splitKeepDelims(combined, []string{TagStart, TagEnd})

	var out []string
	pureLen := 0
	for _, p := range parts {
		if p == "" {
			continue
		}
		if p == TagStart || p == TagEnd {
			out = append(out, p)
			continue
		}
		remain := maxContent - pureLen
		if remain <= 0 {
			break
		}
		runes := []rune(p)
		if len(runes) <= remain {
			out = append(out, p)
			pureLen += len(runes)
		} else {
			out = append(out, string(runes[:remain]))
			pureLen += remain
			break
		}
	}

	res := strings.Join(out, "")
	if strings.Count(res, TagStart) > strings.Count(res, TagEnd) {
		res += TagEnd
	}
	if strings.Count(res, TagEnd) > strings.Count(res, TagStart) {
		res = strings.Replace(res, TagEnd, "", 1)
	}
	return res, true
}

func chunkText(content string, pageSize int, overlap int) []string {
	r := []rune(content)
	if len(r) <= pageSize {
		return []string{content}
	}
	var chunks []string
	start := 0
	for start < len(r) {
		end := start + pageSize
		if end > len(r) {
			end = len(r)
		}
		chunks = append(chunks, string(r[start:end]))
		start += (pageSize - overlap)
		if start >= len(r)-overlap {
			break
		}
	}
	return chunks
}

func prefixRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// ---- job queue ----

type JobType string

const (
	JobUpsert JobType = "upsert"
	JobDelete JobType = "delete"
)

type Job struct {
	Type  JobType
	DocID string
	Title string
	File  []byte
}

type Server struct {
	cfg        Config
	es         *elasticsearch.Client
	httpClient *http.Client
	jobQ       chan Job
}

func NewServer(cfg Config) (*Server, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ESHost},
	})
	if err != nil {
		return nil, err
	}
	return &Server{
		cfg: cfg,
		es:  es,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		jobQ: make(chan Job, cfg.QueueLimit),
	}, nil
}

func (s *Server) StartWorkers() {
	for i := 0; i < s.cfg.MaxWorkers; i++ {
		go func(workerID int) {
			for job := range s.jobQ {
				switch job.Type {
				case JobUpsert:
					s.handleUpsertJob(job)
				case JobDelete:
					s.handleDeleteJob(job)
				default:
					log.Printf("[worker %d] unknown job type: %s", workerID, job.Type)
				}
			}
		}(i)
	}
}

func (s *Server) enqueue(job Job) bool {
	select {
	case s.jobQ <- job:
		return true
	default:
		return false
	}
}

func (s *Server) sendCallback(payload any) {
	if strings.TrimSpace(s.cfg.CallbackURL) == "" {
		return
	}
	go func() {
		b, _ := json.Marshal(payload)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.CallbackURL, bytes.NewReader(b))
		if err != nil {
			log.Printf("[callback] build req err: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.httpClient.Do(req)
		if err != nil {
			log.Printf("[callback] post err: %v", err)
			return
		}
		_ = resp.Body.Close()
	}()
}

func (s *Server) extractTextWithTika(file []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.cfg.TikaURL, bytes.NewReader(file))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("X-Tika-Charset", "UTF-8")
	req.Header.Set("X-Tika-PDFocrStrategy", "no_ocr")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return "", fmt.Errorf("tika status=%d body=%s", resp.StatusCode, string(body))
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func (s *Server) esDeleteByDocID(ctx context.Context, docID string) error {
	body := map[string]any{
		"query": map[string]any{
			"term": map[string]any{"doc_id": docID},
		},
	}
	b, _ := json.Marshal(body)

	res, err := s.es.DeleteByQuery(
		[]string{s.cfg.IndexName},
		bytes.NewReader(b),
		s.es.DeleteByQuery.WithContext(ctx),
		s.es.DeleteByQuery.WithRefresh(true),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		raw, _ := io.ReadAll(io.LimitReader(res.Body, 8<<10))
		return fmt.Errorf("es delete_by_query error: %s", string(raw))
	}
	return nil
}

func (s *Server) esBulkIndexChunks(ctx context.Context, docID, title string, chunks []string) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	for i, c := range chunks {
		action := map[string]any{
			"index": map[string]any{
				"_index": s.cfg.IndexName,
				"_id":    fmt.Sprintf("%s_%d", docID, i),
			},
		}
		if err := enc.Encode(action); err != nil {
			return err
		}
		doc := map[string]any{
			"doc_id":  docID,
			"title":   title,
			"content": c,
		}
		if err := enc.Encode(doc); err != nil {
			return err
		}
	}

	res, err := s.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		s.es.Bulk.WithContext(ctx),
		s.es.Bulk.WithRefresh("wait_for"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		raw, _ := io.ReadAll(io.LimitReader(res.Body, 16<<10))
		return fmt.Errorf("es bulk error: %s", string(raw))
	}

	var parsed map[string]any
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return err
	}
	if errorsVal, ok := parsed["errors"].(bool); ok && errorsVal {
		return fmt.Errorf("es bulk has item errors")
	}
	return nil
}

func (s *Server) handleUpsertJob(job Job) {
	status := "success"
	errMsg := ""
	pages := 0

	defer func() {
		s.sendCallback(map[string]any{
			"type":   "upsert",
			"doc_id": job.DocID,
			"status": status,
			"pages":  pages,
			"error":  errMsg,
		})
	}()

	content, err := s.extractTextWithTika(job.File)
	if err != nil {
		status = "failed"
		errMsg = err.Error()
		return
	}

	chunks := chunkText(content, 20000, 30)
	pages = len(chunks)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := s.esDeleteByDocID(ctx, job.DocID); err != nil {
		status = "failed"
		errMsg = err.Error()
		return
	}
	if err := s.esBulkIndexChunks(ctx, job.DocID, job.Title, chunks); err != nil {
		status = "failed"
		errMsg = err.Error()
		return
	}
}

func (s *Server) handleDeleteJob(job Job) {
	status := "success"
	defer func() {
		s.sendCallback(map[string]any{
			"type":   "delete",
			"doc_id": job.DocID,
			"status": status,
		})
	}()

	realID, err := url.PathUnescape(job.DocID)
	if err != nil {
		realID = job.DocID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := s.esDeleteByDocID(ctx, realID); err != nil {
		status = "failed"
	}
}

// ---- HTTP helpers ----

func (s *Server) UpsertHandler(w http.ResponseWriter, r *http.Request) {
	// 防止大文件撑爆内存（按需调整）
	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)

	if err := r.ParseMultipartForm(200 << 20); err != nil {
		writeJSON(w, 400, map[string]any{"error": "invalid multipart form"})
		return
	}

	docID := r.FormValue("doc_id")
	title := r.FormValue("title")
	file, _, err := r.FormFile("file")
	if docID == "" || title == "" || err != nil {
		writeJSON(w, 400, map[string]any{"error": "Missing parameters"})
		return
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, 400, map[string]any{"error": "read file failed"})
		return
	}

	if !s.enqueue(Job{Type: JobUpsert, DocID: docID, Title: title, File: b}) {
		writeJSON(w, 503, map[string]any{"error": "System queue is full"})
		return
	}
	writeJSON(w, 202, map[string]any{"status": "accepted", "doc_id": docID})
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	docID := strings.TrimPrefix(r.URL.Path, "/delete/")
	if docID == "" {
		writeJSON(w, 400, map[string]any{"error": "Missing doc_id"})
		return
	}
	if !s.enqueue(Job{Type: JobDelete, DocID: docID}) {
		writeJSON(w, 503, map[string]any{"error": "System busy"})
		return
	}
	writeJSON(w, 202, map[string]any{"status": "accepted", "doc_id": docID})
}

func intFromQuery(r *http.Request, key string, def int) int {
	v := strings.TrimSpace(r.URL.Query().Get(key))
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

// ✅ 补丁后的 search：collapse + inner_hits(best_chunk)，取每个 doc 的最强 chunk
func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	scope := strings.ToLower(r.URL.Query().Get("scope"))
	if scope == "" {
		scope = "all"
	}

	topFragments := intFromQuery(r, "top_fragments", s.cfg.DefaultTopFragments)
	maxContent := intFromQuery(r, "max_content", s.cfg.DefaultMaxContent)
	msm := r.URL.Query().Get("msm")
	if msm == "" {
		msm = s.cfg.DefaultMSM
	}

	page := intFromQuery(r, "page", 1)
	if page < 1 {
		page = 1
	}
	pageSize := intFromQuery(r, "page_size", s.cfg.DefaultPageSize)
	if pageSize < 1 {
		pageSize = s.cfg.DefaultPageSize
	}
	if pageSize > s.cfg.MaxPageSize {
		pageSize = s.cfg.MaxPageSize
	}
	from := (page - 1) * pageSize
	if from > s.cfg.MaxFrom {
		writeJSON(w, 400, map[string]any{"error": "Page limit exceeded"})
		return
	}

	fields := []string{"title^2", "content"}
	if scope == "title" {
		fields = []string{"title"}
	}

	searchQuery := map[string]any{
		"multi_match": map[string]any{
			"query":                q,
			"fields":               fields,
			"minimum_should_match": msm,
			"analyzer":             "ik_stop_analyzer",
		},
	}

	// ✅ 重点：highlight 放到 inner_hits 里；外层仅负责 doc_id 去重 + aggs 统计
	body := map[string]any{
		"from":  from,
		"size":  pageSize,
		"query": searchQuery,
		"collapse": map[string]any{
			"field": "doc_id",
			"inner_hits": map[string]any{
				"name": "best_chunk",
				"size": 1,
				"sort": []any{
					map[string]any{"_score": "desc"},
				},
				"_source": []string{"doc_id", "title", "content"},
				"highlight": map[string]any{
					"type":                "unified",
					"pre_tags":            []string{TagStart},
					"post_tags":           []string{TagEnd},
					"require_field_match": false,
					"fields": map[string]any{
						"content": map[string]any{
							"number_of_fragments": topFragments,
							"fragment_size":       s.cfg.FragmentSize,
						},
						"title": map[string]any{
							"number_of_fragments": 0,
						},
					},
				},
			},
		},
		"aggs": map[string]any{
			"unique_doc_count": map[string]any{
				"cardinality": map[string]any{
					"field":               "doc_id",
					"precision_threshold": 3000,
				},
			},
		},
	}

	b, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(s.cfg.IndexName),
		s.es.Search.WithBody(bytes.NewReader(b)),
	)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "ES Error: " + err.Error()})
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		raw, _ := io.ReadAll(io.LimitReader(res.Body, 8<<10))
		writeJSON(w, 500, map[string]any{"error": "ES Error: " + string(raw)})
		return
	}

	var parsed map[string]any
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		writeJSON(w, 500, map[string]any{"error": "ES response decode error"})
		return
	}

	totalDocs := 0
	if aggs, ok := parsed["aggregations"].(map[string]any); ok {
		if u, ok := aggs["unique_doc_count"].(map[string]any); ok {
			if v, ok := u["value"].(float64); ok {
				totalDocs = int(v)
			}
		}
	}

	items := make([]map[string]any, 0)
	hitsObj, _ := parsed["hits"].(map[string]any)
	hitsArr, _ := hitsObj["hits"].([]any)

	for _, h := range hitsArr {
		outerHit, _ := h.(map[string]any)

		// ✅ 从 inner_hits.best_chunk.hits.hits[0] 取“最强 chunk”
		innerHits, _ := outerHit["inner_hits"].(map[string]any)
		bestChunkAny, ok := innerHits["best_chunk"].(map[string]any)
		if !ok {
			continue
		}
		bestHitsObj, _ := bestChunkAny["hits"].(map[string]any)
		bestHitsArr, _ := bestHitsObj["hits"].([]any)
		if len(bestHitsArr) == 0 {
			continue
		}

		bestHit, _ := bestHitsArr[0].(map[string]any)
		bestSrc, _ := bestHit["_source"].(map[string]any)
		bestHL, _ := bestHit["highlight"].(map[string]any)

		docID, _ := bestSrc["doc_id"].(string)
		title, _ := bestSrc["title"].(string)
		content, _ := bestSrc["content"].(string)

		// title highlight
		hTitle := mergeHighlightTags(title)
		if hlTitleAny, ok := bestHL["title"]; ok {
			if hlTitleArr, ok := hlTitleAny.([]any); ok && len(hlTitleArr) > 0 {
				if s0, ok := hlTitleArr[0].(string); ok {
					hTitle = mergeHighlightTags(s0)
				}
			}
		}

		// content highlight fragments
		rawFragments := []string{}
		if hlContAny, ok := bestHL["content"]; ok {
			if hlContArr, ok := hlContAny.([]any); ok {
				for _, x := range hlContArr {
					if sx, ok := x.(string); ok {
						rawFragments = append(rawFragments, sx)
					}
				}
			}
		}

		finalContent := ""
		hitSeqCount := 0
		if hBody, ok := processHighlightContent(rawFragments, maxContent); ok {
			for _, f := range rawFragments {
				if strings.Contains(f, TagStart) {
					hitSeqCount++
				}
			}
			finalContent = "..." + strings.TrimSpace(hBody) + "..."
		} else {
			full := sanitizeText(content)
			finalContent = prefixRunes(full, maxContent) + "..."
			hitSeqCount = 0
		}

		items = append(items, map[string]any{
			"doc_id":            docID,
			"title":             title,
			"highlight_title":   hTitle,
			"highlight_content": finalContent,
			"hit_seq_count":     hitSeqCount,
		})
	}

	writeJSON(w, 200, map[string]any{
		"total":       totalDocs,
		"page":        page,
		"page_size":   pageSize,
		"items_count": len(items),
		"items":       items,
	})
}

// ---- (可选) 如果你以后想发 multipart 请求，可用这个 helper，当前未使用 ----
func _buildMultipart(docID, title string, fileField string, fileName string, fileContent []byte) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("doc_id", docID)
	_ = w.WriteField("title", title)

	fw, err := w.CreateFormFile(fileField, fileName)
	if err != nil {
		return nil, "", err
	}
	if _, err := fw.Write(fileContent); err != nil {
		return nil, "", err
	}
	_ = w.Close()
	return &buf, w.FormDataContentType(), nil
}
