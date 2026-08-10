package notifier

import (
	"context"
	"errors"
	"testing"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestNotifier_Notify_GetTokenError(t *testing.T) {
	t.Parallel()

	getErr := errors.New("dynamodb unavailable")

	tokenStore := &stubTokenStore{
		err: getErr,
	}
	client := &stubMessageClient{}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		context.Background(),
		"T123",
		"test message",
	)
	if err == nil {
		t.Fatal("Notify() error = nil, want an error")
	}

	if !errors.Is(err, getErr) {
		t.Errorf(
			"Notify() error = %v, want wrapped error %v",
			err,
			getErr,
		)
	}

	const wantError = "get slack token: dynamodb unavailable"
	if err.Error() != wantError {
		t.Errorf(
			"Notify() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if client.called {
		t.Error("PostMessage() was called after Get() failed")
	}
}

func TestNotifier_Notify_PostMessageError(t *testing.T) {
	t.Parallel()

	postErr := errors.New("slack unavailable")

	tokenStore := &stubTokenStore{
		token: &oauth.Token{
			TeamID:      "T123",
			AccessToken: "xoxb-test",
			BotUserID:   "B123",
			ChannelID:   "C123",
		},
	}
	client := &stubMessageClient{
		err: postErr,
	}

	notifier := NewNotifier(tokenStore, client)

	err := notifier.Notify(
		context.Background(),
		"T123",
		"test message",
	)
	if err == nil {
		t.Fatal("Notify() error = nil, want an error")
	}

	if !errors.Is(err, postErr) {
		t.Errorf(
			"Notify() error = %v, want wrapped error %v",
			err,
			postErr,
		)
	}

	const wantError = "post slack message: slack unavailable"
	if err.Error() != wantError {
		t.Errorf(
			"Notify() error = %q, want %q",
			err.Error(),
			wantError,
		)
	}

	if !tokenStore.called {
		t.Fatal("Get() was not called")
	}

	if !client.called {
		t.Fatal("PostMessage() was not called")
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
