package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_PostMessage_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"method = %q, want %q",
					r.Method,
					http.MethodPost,
				)
				return
			}

			wantAuthorization := "Bearer " + testAccessToken
			if got := r.Header.Get("Authorization"); got != wantAuthorization {
				t.Errorf(
					"Authorization = %q, want %q",
					got,
					wantAuthorization,
				)
			}

			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf(
					"Content-Type = %q, want %q",
					got,
					"application/json",
				)
			}

			var request postMessageRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Errorf("failed to decode request body: %v", err)
				return
			}

			if request.Channel != testChannelID {
				t.Errorf(
					"Channel = %q, want %q",
					request.Channel,
					testChannelID,
				)
			}

			if request.Text != testMessageText {
				t.Errorf(
					"Text = %q, want %q",
					request.Text,
					testMessageText,
				)
			}

			w.Header().Set("Content-Type", "application/json")

			response := postMessageResponse{
				OK: true,
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	err := client.PostMessage(
		context.Background(),
		testAccessToken,
		testChannelID,
		testMessageText,
	)
	if err != nil {
		t.Fatalf("PostMessage() error = %v", err)
	}
}

func TestNewClient_UsesDefaultHTTPClient(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	if client.httpClient != http.DefaultClient {
		t.Errorf(
			"httpClient = %p, want http.DefaultClient %p",
			client.httpClient,
			http.DefaultClient,
		)
	}

	if client.postMessageURL != slackPostMessageURL {
		t.Errorf(
			"postMessageURL = %q, want %q",
			client.postMessageURL,
			slackPostMessageURL,
		)
	}
}
