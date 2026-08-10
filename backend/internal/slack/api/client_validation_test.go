package api

import (
	"context"
	"strings"
	"testing"
)

func TestClient_PostMessage_EmptyAccessToken(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	err := client.PostMessage(
		context.Background(),
		"",
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
		"access token is empty",
	) {
		t.Errorf(
			"PostMessage() error = %q, want access token error",
			err,
		)
	}
}

func TestClient_PostMessage_EmptyChannelID(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	err := client.PostMessage(
		context.Background(),
		testAccessToken,
		"",
		testMessage,
	)

	if err == nil {
		t.Fatal(
			"PostMessage() error = nil, want an error",
		)
	}

	if !strings.Contains(
		err.Error(),
		"channel ID is empty",
	) {
		t.Errorf(
			"PostMessage() error = %q, want channel ID error",
			err,
		)
	}
}

func TestClient_PostMessage_EmptyText(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	err := client.PostMessage(
		context.Background(),
		testAccessToken,
		testChannelID,
		"",
	)

	if err == nil {
		t.Fatal(
			"PostMessage() error = nil, want an error",
		)
	}

	if !strings.Contains(
		err.Error(),
		"text is empty",
	) {
		t.Errorf(
			"PostMessage() error = %q, want text error",
			err,
		)
	}
}

func TestClient_OpenConversation_EmptyAccessToken(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	channelID, err := client.OpenConversation(
		context.Background(),
		"",
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
		"access token is empty",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want access token error",
			err,
		)
	}
}

func TestClient_OpenConversation_EmptyUserID(t *testing.T) {
	t.Parallel()

	client := NewClient(nil)

	channelID, err := client.OpenConversation(
		context.Background(),
		testAccessToken,
		"",
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
		"user ID is empty",
	) {
		t.Errorf(
			"OpenConversation() error = %q, want user ID error",
			err,
		)
	}
}
