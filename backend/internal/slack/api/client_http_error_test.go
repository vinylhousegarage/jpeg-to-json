package api

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PostMessage_HTTPError(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusInternalServerError)
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

	if err == nil {
		t.Fatal(
			"PostMessage() error = nil, want an error",
		)
	}

	if !strings.Contains(
		err.Error(),
		"returned status 500",
	) {
		t.Errorf(
			"PostMessage() error = %q, want status error",
			err,
		)
	}
}

func TestClient_PostMessage_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			if _, err := w.Write([]byte("{")); err != nil {
				t.Errorf(
					"failed to write response: %v",
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

	if err == nil {
		t.Fatal(
			"PostMessage() error = nil, want an error",
		)
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

func TestClient_OpenConversation_HTTPError(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	channelID, err := client.OpenConversation(
		context.Background(),
		testAccessToken,
		testUserID,
	)

	if err == nil {
		t.Fatal(
			"OpenConversation() error = nil, want an error",
		)
	}

	if channelID != "" {
		t.Errorf(
			"OpenConversation() channelID = %q, want empty",
			channelID,
		)
	}

	if !strings.Contains(
		err.Error(),
		"returned status 500",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want status error",
			err,
		)
	}
}

func TestClient_OpenConversation_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := newTestServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			_ *http.Request,
		) {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			if _, err := w.Write([]byte("{")); err != nil {
				t.Errorf(
					"failed to write response: %v",
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

	if err == nil {
		t.Fatal(
			"OpenConversation() error = nil, want an error",
		)
	}

	if channelID != "" {
		t.Errorf(
			"OpenConversation() channelID = %q, want empty",
			channelID,
		)
	}

	if !strings.Contains(
		err.Error(),
		"decode slack open conversation response",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want decode error",
			err,
		)
	}
}
