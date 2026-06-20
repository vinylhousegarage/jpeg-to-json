package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

func TestEnableCORS(t *testing.T) {
	// 期待するオリジン（テスト用）
	origin := "http://localhost:3000"
    
	// OPTIONSリクエストテスト
	t.Run("OPTIONS request returns 200 and headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		w := httptest.NewRecorder()

		// 関数呼び出し
		enableCORS(w, req, origin)

		// ステータスコード検証
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		// ヘッダー検証
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("expected %s, got %s", origin, got)
		}
	})

	// POSTリクエストテスト
	t.Run("POST request sets headers", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()

		enableCORS(w, req, origin)

		// ヘッダーがセットされているか確認
		if got := w.Header().Get("Access-Control-Allow-Methods"); got != "POST, OPTIONS" {
			t.Errorf("expected POST, OPTIONS, got %s", got)
		}
	})
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
