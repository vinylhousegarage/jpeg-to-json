package presign

import (
	"bytes"
 "net/http"
 "net/http/httptest"
 "testing"

 	"go.uber.org/zap"
)

func TestPresignHandler_ServeHTTP_Integration(t *testing.T) {
	t.Parallel()

	client := NewPresignClient(&mockPresigner{}, "test-bucket")
	h := NewPresignHandler(
		[]string{"http://localhost:3000"},
		client,
		zap.NewNop(),
	)

	t.Run("CORS: Valid Origin", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/presign", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Error("CORS header should be set for valid origin")
		}
	})

	t.Run("CORS: Invalid Origin (Denied)", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/presign", nil)
		req.Header.Set("Origin", "http://hacker.com")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Error("CORS header should NOT be set for invalid origin")
		}
	})

	t.Run("OPTIONS: Preflight Request", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodOptions, "/presign", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for OPTIONS, got %d", w.Code)
		}
	})

	t.Run("Validation: Invalid Request (400 Bad Request)", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/presign", bytes.NewBufferString(`{"invalid":`))
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
		}
	})

	t.Run("Success: Full Flow (200 OK)", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/presign", bytes.NewBufferString(`{"filename":"test.jpg"}`))
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("upload_url")) {
			t.Error("response body does not contain upload_url")
		}
	})
}
