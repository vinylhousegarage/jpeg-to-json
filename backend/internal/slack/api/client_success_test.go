package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestClient_PostMessage_Success(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"method = %q, want %q",
					r.Method,
					http.MethodPost,
				)
			}

			if got := r.Header.Get("Authorization"); got != "Bearer "+testAccessToken {
				t.Errorf(
					"Authorization = %q, want %q",
					got,
					"Bearer "+testAccessToken,
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
				t.Errorf(
					"failed to decode request: %v",
					err,
				)
				return
			}

			if request.Channel != testChannelID {
				t.Errorf(
					"Channel = %q, want %q",
					request.Channel,
					testChannelID,
				)
			}

			if request.Text != testMessage {
				t.Errorf(
					"Text = %q, want %q",
					request.Text,
					testMessage,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			response := postMessageResponse{
				OK: true,
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf(
					"failed to encode response: %v",
					err,
				)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	err := client.PostMessage(
		context.Background(),
		testAccessToken,
		testChannelID,
		testMessage,
	)
	if err != nil {
		t.Fatalf(
			"PostMessage() error = %v",
			err,
		)
	}
}

func TestClient_OpenConversation_Success(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"method = %q, want %q",
					r.Method,
					http.MethodPost,
				)
			}

			if got := r.Header.Get("Authorization"); got != "Bearer "+testAccessToken {
				t.Errorf(
					"Authorization = %q, want %q",
					got,
					"Bearer "+testAccessToken,
				)
			}

			if got := r.Header.Get("Content-Type"); got != "application/json" {
				t.Errorf(
					"Content-Type = %q, want %q",
					got,
					"application/json",
				)
			}

			var request openConversationRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Errorf(
					"failed to decode request: %v",
					err,
				)
				return
			}

			if request.Users != testUserID {
				t.Errorf(
					"Users = %q, want %q",
					request.Users,
					testUserID,
				)
			}

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			response := openConversationResponse{
				OK: true,
			}
			response.Channel.ID = testChannelID

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf(
					"failed to encode response: %v",
					err,
				)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	channelID, err := client.OpenConversation(
		context.Background(),
		testAccessToken,
		testUserID,
	)
	if err != nil {
		t.Fatalf(
			"OpenConversation() error = %v",
			err,
		)
	}

	if channelID != testChannelID {
		t.Errorf(
			"OpenConversation() channelID = %q, want %q",
			channelID,
			testChannelID,
		)
	}
}
