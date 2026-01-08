package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golib/libs/doc_search"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	ESHost        string
	IndexName     string
	Timeout       time.Duration
	DeleteIfExist bool
}

func LoadConfig() Config {
	return Config{
		ESHost:        doc_search.GetEnv("ES_HOST", "http://192.168.12.49:9200"),
		IndexName:     doc_search.GetEnv("INDEX_NAME", "test_hikb"),
		Timeout:       600 * time.Second,
		DeleteIfExist: true,
	}
}

type IndexManager struct {
	es *elasticsearch.Client
}

func NewIndexManager(cfg Config) (*IndexManager, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.ESHost},
		Transport: &http.Transport{
			ResponseHeaderTimeout: cfg.Timeout,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("创建 ES 客户端失败: %w", err)
	}
	return &IndexManager{es: es}, nil
}

func (m *IndexManager) Ping(ctx context.Context) error {
	resp, err := m.es.Ping(m.es.Ping.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}

func (m *IndexManager) IndexExists(ctx context.Context, index string) (bool, error) {
	resp, err := m.es.Indices.Exists([]string{index}, m.es.Indices.Exists.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("检查索引是否存在返回异常状态码: %d", resp.StatusCode)
	}
}

func (m *IndexManager) DeleteIndex(ctx context.Context, index string) error {
	resp, err := m.es.Indices.Delete([]string{index}, m.es.Indices.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("删除索引失败: %s", strings.TrimSpace(string(b)))
	}
	return nil
}

func (m *IndexManager) CreateIndex(ctx context.Context, index string, mapping any) error {
	payload, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("序列化 mapping 失败: %w", err)
	}

	resp, err := m.es.Indices.Create(
		index,
		m.es.Indices.Create.WithContext(ctx),
		m.es.Indices.Create.WithBody(bytes.NewReader(payload)),
	)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.IsError() {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("创建索引失败: %s", strings.TrimSpace(string(b)))
	}
	return nil
}

func BuildIKStopMapping(stopwords []string) map[string]any {
	return map[string]any{
		"settings": map[string]any{
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"ik_stop_analyzer": map[string]any{
						"type":      "custom",
						"tokenizer": "ik_max_word",
						"filter":    []any{"my_stop_filter"},
					},
				},
				"filter": map[string]any{
					"my_stop_filter": map[string]any{
						"type":      "stop",
						"stopwords": stopwords,
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"doc_id": map[string]any{
					"type": "keyword",
				},
				"title": map[string]any{
					"type":            "text",
					"analyzer":        "ik_stop_analyzer",
					"search_analyzer": "ik_smart",
				},
				"content": map[string]any{
					"type":            "text",
					"analyzer":        "ik_stop_analyzer",
					"search_analyzer": "ik_smart",
				},
			},
		},
	}
}

func CoreStopwords() []string {
	return []string{
		// 极高频虚词
		"的", "了", "在", "是", "我", "有", "和", "就", "不", "人", "都", "一", "一个", "上", "也", "很", "到", "说", "要", "去", "你", "会", "着", "没有", "看", "好", "自己", "这", "那", "个", "中", "为", "以", "地", "于", "之", "后", "与", "及", "而", "等",
		// 标点符号（Tika提取常残留）
		"。", "，", "、", "？", "！", "“", "”", "；", "：", "（", "）", "《", "》", "—", "-", ".", "/", "\\", "[", "]", "{", "}", "=", "+", "*", "&", "^", "%", "$", "#", "@", "!", "~", "`", "|",
		// 格式占位符
		"nbsp", "quot", "amp", "lt", "gt",
	}
}

func EnsureIndex(ctx context.Context, mgr *IndexManager, cfg Config, mapping any) error {
	exists, err := mgr.IndexExists(ctx, cfg.IndexName)
	if err != nil {
		return fmt.Errorf("检查索引是否存在失败: %w", err)
	}

	if exists && cfg.DeleteIfExist {
		if err := mgr.DeleteIndex(ctx, cfg.IndexName); err != nil {
			return err
		}
	}

	// 如果存在且不删除，就直接返回（避免创建报错）
	if exists && !cfg.DeleteIfExist {
		return errors.New("索引已存在且 DeleteIfExist=false，未执行创建")
	}

	return mgr.CreateIndex(ctx, cfg.IndexName, mapping)
}

func main() {
	cfg := LoadConfig()
	ctx := context.Background()

	mgr, err := NewIndexManager(cfg)
	if err != nil {
		panic(err)
	}

	if err := mgr.Ping(ctx); err != nil {
		fmt.Println("连接失败，请检查防火墙设置。错误:", err)
		return
	}
	fmt.Println("成功连接到局域网 ES 服务！")

	mapping := BuildIKStopMapping(CoreStopwords())
	if err := EnsureIndex(ctx, mgr, cfg, mapping); err != nil {
		panic(err)
	}

	fmt.Println("使用 IK 分词器的分段索引构建成功！")
}
