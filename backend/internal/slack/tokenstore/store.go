package tokenstore

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
