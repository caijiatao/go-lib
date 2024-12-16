package http_helper

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func Post(ctx context.Context, url, contentType string, params interface{}, result interface{}) (err error) {
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(paramsStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = HandleResponse(ctx, resp, result)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func HandleResponse(ctx context.Context, resp *http.Response, result interface{}) error {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return HttpResponseCodeError(resp.StatusCode)
	}
	// 确保在body内容读取完成后正确关闭
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, result); err != nil {
		return err
	}
	return nil
}
