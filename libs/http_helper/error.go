package http_helper

import (
	"errors"
	"fmt"
)

var (
	HttpResponseCodeError = func(statusCode int) error {
		return errors.New(fmt.Sprintf("http response error :%d", statusCode))
	}
)
