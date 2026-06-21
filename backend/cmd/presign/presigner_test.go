package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

func TestSetCORSHeaders(t *testing.T) {
	origin := "http://localhost:3000"
	w := httptest.NewRecorder()

	setCORSHeaders(w, origin)

	expected := map[string]string{
		"Access-Control-Allow-Origin":  origin,
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type",
	}

	for header, want := range expected {
		got := w.Header().Get(header)
		if got != want {
			t.Errorf("expected header %s to be %q, got %q", header, want, got)
		}
	}
}

const StatusCodeIgnore = 0

func TestValidateRequest(t *testing.T) {
	t.Parallel()
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
			got, err := validateRequest(req)

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

// テスト用モック
type mockPresigner struct{}

func (m *mockPresigner) PresignPutObject(
  ctx context.Context,
  params *s3.PutObjectInput,
  optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	return &v4.PresignedHTTPRequest{URL: "https://example.com/test.jpg"}, nil
}

func TestGeneratePresignURL(t *testing.T) {
	// コンストラクタを使用
	client := newPresignClient(&mockPresigner{}, "test-bucket")

	url, expiresAt, err := client.generatePresignURL(context.Background(), "test.jpg")

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
	// テスト用データ
	url := "https://example.com/upload"
	duration := 15 * time.Minute
	expiresAt := time.Now().Add(duration).Truncate(time.Second)

	// HTTPレスポンスレコーダー
	w := httptest.NewRecorder()

	// 関数を実行
	writeJSON(w, url, expiresAt)

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
