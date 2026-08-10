package api

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestClient_PostMessage_SlackError(t *testing.T) {
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

			if _, err := w.Write(
				[]byte(`{"ok":false,"error":"channel_not_found"}`),
			); err != nil {
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
		"channel_not_found",
	) {
		t.Errorf(
			"PostMessage() error = %q, want Slack API error",
			err,
		)
	}
}

func TestClient_OpenConversation_SlackError(t *testing.T) {
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

			if _, err := w.Write(
				[]byte(`{"ok":false,"error":"user_not_found"}`),
			); err != nil {
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
		"user_not_found",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want Slack API error",
			err,
		)
	}
}

func TestClient_OpenConversation_MissingChannelID(t *testing.T) {
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

			if _, err := w.Write(
				[]byte(`{"ok":true,"channel":{}}`),
			); err != nil {
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
		"missing channel.id",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want missing channel.id error",
			err,
		)
	}
}
