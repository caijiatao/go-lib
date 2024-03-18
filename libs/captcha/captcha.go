package captcha

import (
	"bytes"
	"errors"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
	"time"
)

type CaptchaServer struct {
	ID   string //验证码ID
	Code string //验证码数字
}

func NewCaptchaID() string {
	return captcha.New()
}

func NewCaptcha() *CaptchaServer {
	return &CaptchaServer{}
}

func (self *CaptchaServer) CaptchaID(id string) *CaptchaServer {
	self.ID = id
	return self
}

func (self *CaptchaServer) CaptchaCode(code string) *CaptchaServer {
	self.Code = code
	return self
}

func (self CaptchaServer) Verify() error {
	if !captcha.VerifyString(self.ID, self.Code) {
		return errors.New("验证失败")
	}
	return nil
}

func Server(ctx *gin.Context) error {
	server := NewCaptcha()
	err := server.ServeHTTP(ctx.Writer, ctx.Request)
	if err != nil {
		return err
	}
	return nil
}

func (h *CaptchaServer) serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		err := captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)
		if err != nil {
			return err
		}
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		err := captcha.WriteAudio(&content, id, lang)
		if err != nil {
			return err
		}
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

func (h *CaptchaServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		return errors.New("bad request")
	}
	if r.FormValue("reload") == "true" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	err := h.serve(w, r, id, ext, lang, download)
	if err != nil {
		return err
	}
	return nil
}
