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

type stubDynamoDBClient struct {
	putItemOutput *dynamodb.PutItemOutput
	putItemErr    error
	putItemCalled bool
	putItemInput  *dynamodb.PutItemInput

	getItemOutput *dynamodb.GetItemOutput
	getItemErr    error
	getItemCalled bool
	getItemInput  *dynamodb.GetItemInput
}

func (s *stubDynamoDBClient) PutItem(
	_ context.Context,
	input *dynamodb.PutItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	s.putItemCalled = true
	s.putItemInput = input

	return s.putItemOutput, s.putItemErr
}

func (s *stubDynamoDBClient) GetItem(
	_ context.Context,
	input *dynamodb.GetItemInput,
	_ ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	s.getItemCalled = true
	s.getItemInput = input

	return s.getItemOutput, s.getItemErr
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

	client := &stubDynamoDBClient{
		putItemOutput: &dynamodb.PutItemOutput{},
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

	if !client.putItemCalled {
		t.Fatal("PutItem() was not called")
	}

	if client.putItemInput == nil {
		t.Fatal("PutItem() input is nil")
	}

	if client.putItemInput.TableName == nil {
		t.Fatal("PutItem() TableName is nil")
	}

	if got := *client.putItemInput.TableName; got != testTableName {
		t.Errorf(
			"PutItem() TableName = %q, want %q",
			got,
			testTableName,
		)
	}

	var item tokenItem
	if err := attributevalue.UnmarshalMap(
		client.putItemInput.Item,
		&item,
	); err != nil {
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

	if client.getItemCalled {
		t.Error("GetItem() was called by Save()")
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

			client := &stubDynamoDBClient{}
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

			if client.putItemCalled {
				t.Error("PutItem() was called for invalid token")
			}

			if client.getItemCalled {
				t.Error("GetItem() was called by Save()")
			}
		})
	}
}

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
