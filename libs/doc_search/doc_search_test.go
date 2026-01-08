package doc_search

import (
	"testing"
)

func TestMain(m *testing.M) {

	m.Run()
}

func Test(t *testing.T) {

}

func TestServer_Analyzer(t *testing.T) {

	cfg := LoadConfig()
	server, _ := NewServer(cfg)
	tests := []struct {
		name    string
		reqBody map[string]any
	}{

		{
			name: "空格分词",
			reqBody: map[string]any{
				"analyzer": "whitespace",
				"text":     "2 running Quick brown-foxes leap over lazy dogs in the summer evening.",
			},
		},
		{
			name: "标准分词",
			reqBody: map[string]any{
				"analyzer": "standard",
				"text":     "2 running Quick brown-foxes leap over lazy dogs in the summer evening.",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server.Analyzer(test.reqBody)
		})
	}
}
