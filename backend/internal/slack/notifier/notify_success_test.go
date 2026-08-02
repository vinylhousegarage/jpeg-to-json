package notifier

import (
	"context"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestNotifier_Notify_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tokenStore := &stubTokenStore{
		token: &oauth.Token{
			TeamID:      "T123",
			AccessToken: "xoxb-test",
			BotUserID:   "B123",
		},
	}
	client := &stubMessageClient{}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		ctx,
		"T123",
		"C123",
		"test message",
	)
	if err != nil {
		t.Fatalf("Notify() error = %v", err)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if tokenStore.ctx != ctx {
		t.Error("Get() received an unexpected context")
	}

	if tokenStore.teamID != "T123" {
		t.Errorf(
			"Get() teamID = %q, want %q",
			tokenStore.teamID,
			"T123",
		)
	}

	if !client.called {
		t.Fatal("PostMessage() was not called")
	}

	if client.ctx != ctx {
		t.Error("PostMessage() received an unexpected context")
	}

	if client.accessToken != "xoxb-test" {
		t.Errorf(
			"PostMessage() accessToken = %q, want %q",
			client.accessToken,
			"xoxb-test",
		)
	}

	if client.channelID != "C123" {
		t.Errorf(
			"PostMessage() channelID = %q, want %q",
			client.channelID,
			"C123",
		)
	}

	if client.text != "test message" {
		t.Errorf(
			"PostMessage() text = %q, want %q",
			client.text,
			"test message",
		)
	}
}
