package storage

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
		client,
		zap.NewNop(),
	)

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
		if !bytes.Contains(w.Body.Bytes(), []byte("uploadURL")) {
			t.Error("response body does not contain uploadURL")
		}
	})
}
