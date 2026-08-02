package tokenstore

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

type putItemAPI interface {
	PutItem(
		ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// Store persists Slack OAuth tokens in DynamoDB.
type Store struct {
	client    putItemAPI
	tableName string
	now       func() time.Time
}

func NewStore(
	client putItemAPI,
	tableName string,
) *Store {
	return &Store{
		client:    client,
		tableName: tableName,
		now:       time.Now,
	}
}

// tokenItem represents a Slack bot token persisted in DynamoDB.
type tokenItem struct {
	TeamID      string `dynamodbav:"team_id"`
	AccessToken string `dynamodbav:"access_token"`
	BotUserID   string `dynamodbav:"bot_user_id"`
	UpdatedAt   string `dynamodbav:"updated_at"`
}

func (s *Store) Save(
	ctx context.Context,
	token *oauth.Token,
) error {
	if token == nil {
		return fmt.Errorf("save slack token: token is nil")
	}

	if token.TeamID == "" {
		return fmt.Errorf("save slack token: team ID is empty")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("save slack token: access token is empty")
	}

	if token.BotUserID == "" {
		return fmt.Errorf("save slack token: bot user ID is empty")
	}

	item := tokenItem{
		TeamID:      token.TeamID,
		AccessToken: token.AccessToken,
		BotUserID:   token.BotUserID,
		UpdatedAt:   s.now().UTC().Format(time.RFC3339),
	}

	attributes, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal slack token item: %w", err)
	}

	_, err = s.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: &s.tableName,
			Item:      attributes,
		},
	)
	if err != nil {
		return fmt.Errorf("put slack token item: %w", err)
	}

	return nil
}
