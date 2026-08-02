package tokenstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

var ErrTokenNotFound = errors.New("slack token not found")

func (s *Store) Get(
	ctx context.Context,
	teamID string,
) (*oauth.Token, error) {
	if teamID == "" {
		return nil, fmt.Errorf("get slack token: team ID is empty")
	}

	output, err := s.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			TableName: &s.tableName,
			Key: map[string]types.AttributeValue{
				"team_id": &types.AttributeValueMemberS{
					Value: teamID,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get slack token item: %w", err)
	}

	if output == nil || len(output.Item) == 0 {
		return nil, fmt.Errorf(
			"get slack token for team %q: %w",
			teamID,
			ErrTokenNotFound,
		)
	}

	var item tokenItem
	if err := attributevalue.UnmarshalMap(output.Item, &item); err != nil {
		return nil, fmt.Errorf("unmarshal slack token item: %w", err)
	}

	if item.TeamID == "" {
		return nil, fmt.Errorf("get slack token: stored team ID is empty")
	}

	if item.AccessToken == "" {
		return nil, fmt.Errorf("get slack token: stored access token is empty")
	}

	if item.BotUserID == "" {
		return nil, fmt.Errorf("get slack token: stored bot user ID is empty")
	}

	return &oauth.Token{
		TeamID:      item.TeamID,
		AccessToken: item.AccessToken,
		BotUserID:   item.BotUserID,
	}, nil
}
