package bedrock

import (
	"testing"
)

func TestParseResponse(t *testing.T) {
	s := &Service{client: nil}

	tests := []struct {
		name    string
		input   string
		wantKey string
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			input:   `{"content": [{"text": "{\"key\": \"value\"}"}]}`,
			wantKey: "value",
			wantErr: false,
		},
		{
			name:    "JSON with Markdown",
			input:   "{\"content\": [{\"text\": \"```json\\n{\\\"key\\\": \\\"value\\\"}\\n```\"}]}",
			wantKey: "value",
			wantErr: false,
		},
		{
			name:    "JSON with explanatory text",
			input:   `{"content": [{"text": "解析結果は以下の通りです。{\"key\": \"value\"} 以上です。"}]}`,
			wantKey: "value",
			wantErr: false,
		},
		{
			name:    "Invalid JSON Format",
			input:   `{"content": [{"text": "これはJSONではありません"}]}`,
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Empty Content",
			input:   `{"content": []}`,
			wantKey: "",
			wantErr: true,
		},
		{
			name:    "Broken JSON inside braces",
			input:   `{"content": [{"text": "{\"key\": \"value\""}]}`,
			wantKey: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := s.parseResponse([]byte(tt.input))
			
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				val, ok := res["key"]
				if !ok {
					t.Errorf("expected key 'key' not found in result")
				} else if val != tt.wantKey {
					t.Errorf("got %v, want %v", val, tt.wantKey)
				}
			}
		})
	}
}
