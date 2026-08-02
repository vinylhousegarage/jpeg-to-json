package tokenstore

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestStore_Save_PutItemError(t *testing.T) {
	t.Parallel()

	putErr := errors.New("dynamodb unavailable")
	client := &stubDynamoDBClient{
		putItemErr: putErr,
	}

	store := NewStore(client, testTableName)
	store.now = func() time.Time {
		return time.Date(
			2026,
			time.August,
			2,
			12,
			0,
			0,
			0,
			time.UTC,
		)
	}

	token := &oauth.Token{
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
	}

	err := store.Save(context.Background(), token)
	if err == nil {
		t.Fatal("Save() error = nil, want an error")
	}

	if !errors.Is(err, putErr) {
		t.Errorf(
			"Save() error = %v, want wrapped error %v",
			err,
			putErr,
		)
	}

	const wantError = "put slack token item: dynamodb unavailable"
	if got := err.Error(); got != wantError {
		t.Errorf(
			"Save() error = %q, want %q",
			got,
			wantError,
		)
	}

	if !client.putItemCalled {
		t.Fatal("PutItem() was not called")
	}

	if client.getItemCalled {
		t.Error("GetItem() was called by Save()")
	}
}
