package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestWriteError(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	t.Run("should return the specified status and code when APIError is provided", func(t *testing.T) {
		t.Parallel()

		err := New(ErrorCodeMissingFilename, http.StatusBadRequest, nil)
		w := httptest.NewRecorder()

		WriteError(w, err, logger)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, but got %d", http.StatusBadRequest, w.Code)
		}

		var res ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if res.Error != string(ErrorCodeMissingFilename) {
			t.Errorf("expected error code %s, but got %s", ErrorCodeMissingFilename, res.Error)
		}
	})

	t.Run("should return 500 and internal_server_error for unknown errors", func(t *testing.T) {
		t.Parallel()

		err := errors.New("unknown database error")
		w := httptest.NewRecorder()

		WriteError(w, err, logger)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, but got %d", http.StatusInternalServerError, w.Code)
		}

		var res ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if res.Error != string(ErrorCodeInternal) {
			t.Errorf("expected error code %s, but got %s", ErrorCodeInternal, res.Error)
		}
	})
}
