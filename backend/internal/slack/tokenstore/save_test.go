package tokenstore

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

const testTableName = "slack-tokens"

type stubPutItemClient struct {
	output *dynamodb.PutItemOutput
	err    error

	called bool
	input  *dynamodb.PutItemInput
}

func (s *stubPutItemClient) PutItem(
	_ context.Context,
	input *dynamodb.PutItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	s.called = true
	s.input = input

	return s.output, s.err
}

func TestStore_Save_Success(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(
		2026,
		time.August,
		2,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	client := &stubPutItemClient{
		output: &dynamodb.PutItemOutput{},
	}
	store := NewStore(client, testTableName)
	store.now = func() time.Time {
		return fixedTime
	}

	token := &oauth.Token{
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
	}

	err := store.Save(context.Background(), token)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !client.called {
		t.Fatal("PutItem() was not called")
	}

	if client.input == nil {
		t.Fatal("PutItem() input is nil")
	}

	if client.input.TableName == nil {
		t.Fatal("PutItem() TableName is nil")
	}

	if got := *client.input.TableName; got != testTableName {
		t.Errorf(
			"PutItem() TableName = %q, want %q",
			got,
			testTableName,
		)
	}

	var item tokenItem
	if err := attributevalue.UnmarshalMap(client.input.Item, &item); err != nil {
		t.Fatalf("failed to unmarshal PutItem item: %v", err)
	}

	if item.TeamID != token.TeamID {
		t.Errorf(
			"TeamID = %q, want %q",
			item.TeamID,
			token.TeamID,
		)
	}

	if item.AccessToken != token.AccessToken {
		t.Errorf(
			"AccessToken = %q, want %q",
			item.AccessToken,
			token.AccessToken,
		)
	}

	if item.BotUserID != token.BotUserID {
		t.Errorf(
			"BotUserID = %q, want %q",
			item.BotUserID,
			token.BotUserID,
		)
	}

	wantUpdatedAt := fixedTime.Format(time.RFC3339)
	if item.UpdatedAt != wantUpdatedAt {
		t.Errorf(
			"UpdatedAt = %q, want %q",
			item.UpdatedAt,
			wantUpdatedAt,
		)
	}
}

func TestStore_Save_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		token     *oauth.Token
		wantError string
	}{
		{
			name:      "nil token",
			token:     nil,
			wantError: "save slack token: token is nil",
		},
		{
			name: "missing team ID",
			token: &oauth.Token{
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
			},
			wantError: "save slack token: team ID is empty",
		},
		{
			name: "missing access token",
			token: &oauth.Token{
				TeamID:    "T123",
				BotUserID: "B123",
			},
			wantError: "save slack token: access token is empty",
		},
		{
			name: "missing bot user ID",
			token: &oauth.Token{
				TeamID:      "T123",
				AccessToken: "xoxb-test",
			},
			wantError: "save slack token: bot user ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &stubPutItemClient{}
			store := NewStore(client, testTableName)

			err := store.Save(context.Background(), tt.token)
			if err == nil {
				t.Fatal("Save() error = nil, want an error")
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"Save() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}

			if client.called {
				t.Error("PutItem() was called for invalid token")
			}
		})
	}
}

func TestStore_Save_PutItemError(t *testing.T) {
	t.Parallel()

	putErr := errors.New("dynamodb unavailable")
	client := &stubPutItemClient{
		err: putErr,
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

	if got := err.Error(); got != "put slack token item: dynamodb unavailable" {
		t.Errorf(
			"Save() error = %q, want %q",
			got,
			"put slack token item: dynamodb unavailable",
		)
	}

	if !client.called {
		t.Fatal("PutItem() was not called")
	}
}
