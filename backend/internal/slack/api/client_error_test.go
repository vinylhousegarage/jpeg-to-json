package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_PostMessage_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
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
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "returned status 500") {
		t.Errorf(
			"PostMessage() error = %q, want HTTP status error",
			err,
		)
	}
}

func TestClient_PostMessage_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if _, err := w.Write([]byte("{")); err != nil {
				t.Errorf("failed to write response: %v", err)
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
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	if !strings.Contains(
		err.Error(),
		"decode slack post message response",
	) {
		t.Errorf(
			"PostMessage() error = %q, want decode error",
			err,
		)
	}
}

func TestClient_PostMessage_SlackError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			response := postMessageResponse{
				OK:    false,
				Error: "channel_not_found",
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
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	if !strings.Contains(err.Error(), "channel_not_found") {
		t.Errorf(
			"PostMessage() error = %q, want channel_not_found",
			err,
		)
	}
}

func TestClient_PostMessage_TransportError(t *testing.T) {
	t.Parallel()

	transportErr := errors.New("network unavailable")

	httpClient := &http.Client{
		Transport: roundTripFunc(
			func(*http.Request) (*http.Response, error) {
				return nil, transportErr
			},
		),
	}

	client := NewClient(httpClient)

	err := client.PostMessage(
		context.Background(),
		testAccessToken,
		testChannelID,
		testMessageText,
	)
	if err == nil {
		t.Fatal("PostMessage() error = nil, want an error")
	}

	if !errors.Is(err, transportErr) {
		t.Errorf(
			"PostMessage() error = %v, want wrapped error %v",
			err,
			transportErr,
		)
	}
}
