package map_demo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterUpdateMapNil(t *testing.T) {
	tests := []struct {
		name  string
		m     map[string]interface{}
		wantM map[string]interface{}
	}{
		{
			name:  "nil map",
			m:     nil,
			wantM: nil,
		},
		{
			name:  "empty map",
			m:     map[string]interface{}{},
			wantM: map[string]interface{}{},
		},
		{
			name:  "map with elements",
			m:     map[string]interface{}{"a": "1", "b": "2"},
			wantM: map[string]interface{}{"a": "1", "b": "2"},
		},
		{
			name:  "map with nil value",
			m:     map[string]interface{}{"a": "1", "b": "2", "c": ""},
			wantM: map[string]interface{}{"a": "1", "b": "2", "c": ""},
		},
		{
			name:  "map with nil key",
			m:     map[string]interface{}{"a": "1", "b": "2", "": "3"},
			wantM: map[string]interface{}{"a": "1", "b": "2", "": "3"},
		},
		{
			name:  "map with nil key and value",
			m:     map[string]interface{}{"a": "1", "b": "2", "": ""},
			wantM: map[string]interface{}{"a": "1", "b": "2", "": ""},
		},
		{
			name:  "map with nil key and value",
			m:     map[string]interface{}{"": nil},
			wantM: map[string]interface{}{},
		},
		{
			name:  "map with nil key and value",
			m:     map[string]interface{}{"a": "1", "b": "2", "c": "3", "": nil, "d": ""},
			wantM: map[string]interface{}{"a": "1", "b": "2", "c": "3", "d": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterUpdateMapNil(tt.m)
			assert.Equal(t, tt.wantM, result)
		})
	}
}
