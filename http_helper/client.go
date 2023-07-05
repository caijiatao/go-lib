package http_helper

import "github.com/go-resty/resty/v2"

func NewHTTPClient() {
	client := resty.New()

	client.R().Post()
}
