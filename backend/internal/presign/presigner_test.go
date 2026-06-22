package presign

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"

	"go.uber.org/zap"
)

// テスト用ハンドラー
func setupTestHandler() *PresignHandler {
	return &PresignHandler{
		logger: zap.NewNop(),
		client: &PresignClient{
			s3Presigner: &mockPresigner{},
			bucketName:  "test-bucket",
		},
	}
}

func TestSetCORSHeaders(t *testing.T) {
	t.Parallel()

	// 許可リストを定義
	allowed := []string{"http://localhost:3000"}
	
	// 許可リストを渡す
	h := &PresignHandler{
		allowedOrigins: allowed,
		logger:         zap.NewNop(),
	}

	tests := []struct {
		name          string
		origin        string
		wantHeaderSet bool
	}{
		{"Allowed", "http://localhost:3000", true},
		{"Denied", "http://malicious.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.setCORSHeaders(w, tt.origin)

			got := w.Header().Get("Access-Control-Allow-Origin")
			if tt.wantHeaderSet && got != tt.origin {
				t.Errorf("expected header %s, got %s", tt.origin, got)
			}
			if !tt.wantHeaderSet && got != "" {
				t.Errorf("expected no header, got %s", got)
			}
		})
	}
}

const StatusCodeIgnore = 0

func TestValidateRequest(t *testing.T) {
	t.Parallel()

	h := setupTestHandler()

	tests := []struct {
		name           string
		method         string
		body           string
		wantStatusCode int
		wantErrCode    apierror.ErrorCode
	}{
		{
			name:           "Success: Valid request",
			method:         http.MethodPost,
			body:           `{"filename": "test.txt", "filetype": "text/plain"}`,
			wantStatusCode: StatusCodeIgnore,
		},
		{
			name:           "Error: Method not allowed",
			method:         http.MethodGet,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantErrCode:    apierror.ErrorCodeInvalidMethod,
		},
		{
			name:           "Error: Invalid JSON",
			method:         http.MethodPost,
			body:           `{"filename": "test.txt"`,
			wantStatusCode: http.StatusBadRequest,
			wantErrCode:    apierror.ErrorCodeInvalidJSON,
		},
		{
			name:           "Error: Missing filename",
			method:         http.MethodPost,
			body:           `{"filename": ""}`,
			wantStatusCode: http.StatusBadRequest,
			wantErrCode:    apierror.ErrorCodeMissingFilename,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, "/presign", bytes.NewBufferString(tt.body))
			got, err := h.validateRequest(req)

			if tt.wantErrCode != "" {
				// 異常系
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
				apiErr, ok := err.(*apierror.APIError)
				if !ok {
					t.Fatalf("expected error to be *apierror.APIError, got %T", err)
				}
				if apiErr.Code != tt.wantErrCode {
					t.Errorf("error code mismatch: want %s, got %s", tt.wantErrCode, apiErr.Code)
				}

				// ステータスコードが無視設定でなければ検証
				if tt.wantStatusCode != StatusCodeIgnore {
					if apiErr.HTTPStatus != tt.wantStatusCode {
						t.Errorf("HTTP status mismatch: want %d, got %d", tt.wantStatusCode, apiErr.HTTPStatus)
					}
				}
			} else {
				// 正常系の検証
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
				if got == nil || got.Filename == "" {
					t.Error("request was not parsed correctly")
				}
			}
		})
	}
}

func TestGeneratePresignURL(t *testing.T) {
	t.Parallel()

	h := setupTestHandler()

	url, expiresAt, err := h.generatePresignURL(context.Background(), "test.jpg") // メソッド呼び出しに変更

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "https://example.com/test.jpg"
	if url != expected {
		t.Errorf("expected URL %s, got %s", expected, url)
	}

	if expiresAt.Before(time.Now()) {
		t.Errorf("expected future expiration, got %v", expiresAt)
	}
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()

	h := setupTestHandler()
	url := "https://example.com/upload"
	expiresAt := time.Now().Add(15 * time.Minute).Truncate(time.Second)

	w := httptest.NewRecorder()

	h.writeJSON(w, url, expiresAt)

	// 検証
	resp := w.Result()
	defer resp.Body.Close()

	// ステータスコードのチェック
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", resp.StatusCode)
	}

	// Content-Typeのチェック
	if contentType := resp.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %v", contentType)
	}

	// JSONのデコードと検証
	var got PresignResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.UploadURL != url {
		t.Errorf("expected URL %s, got %s", url, got.UploadURL)
	}

	if got.ExpiresAt.Sub(expiresAt).Abs() > time.Second {
    t.Errorf("expected time near %v, got %v", expiresAt, got.ExpiresAt)
	}
}
