package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendNotification(t *testing.T) {
	t.Parallel()

	const (
		testShotNumber  = "1"
		testDownloadURL = "https://example.com/download"
	)

	expectedText := fmt.Sprintf("撮影番号: %s\n%s", testShotNumber, testDownloadURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var payload SlackPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if payload.Text != expectedText {
			t.Errorf("expected text %q, got %q", expectedText, payload.Text)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(testShotNumber, server.URL)

	err := client.SendNotification(context.Background(), testDownloadURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
