package tokenstore

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestStore_Get_Success(t *testing.T) {
	t.Parallel()

	wantItem := tokenItem{
		TeamID:      "T123",
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		UpdatedAt:   "2026-08-02T12:00:00Z",
	}

	attributes, err := attributevalue.MarshalMap(wantItem)
	if err != nil {
		t.Fatalf("failed to marshal test token item: %v", err)
	}

	client := &stubDynamoDBClient{
		getItemOutput: &dynamodb.GetItemOutput{
			Item: attributes,
		},
	}
	store := NewStore(client, testTableName)

	got, err := store.Get(context.Background(), wantItem.TeamID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if !client.getItemCalled {
		t.Fatal("GetItem() was not called")
	}

	if client.getItemInput == nil {
		t.Fatal("GetItem() input is nil")
	}

	if client.getItemInput.TableName == nil {
		t.Fatal("GetItem() TableName is nil")
	}

	if gotTableName := *client.getItemInput.TableName; gotTableName != testTableName {
		t.Errorf(
			"GetItem() TableName = %q, want %q",
			gotTableName,
			testTableName,
		)
	}

	teamIDAttribute, ok := client.getItemInput.Key["team_id"]
	if !ok {
		t.Fatal(`GetItem() key does not contain "team_id"`)
	}

	teamID, ok := teamIDAttribute.(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf(
			`GetItem() key "team_id" type = %T, want *types.AttributeValueMemberS`,
			teamIDAttribute,
		)
	}

	if teamID.Value != wantItem.TeamID {
		t.Errorf(
			`GetItem() key "team_id" = %q, want %q`,
			teamID.Value,
			wantItem.TeamID,
		)
	}

	if got == nil {
		t.Fatal("Get() token is nil")
	}

	if got.TeamID != wantItem.TeamID {
		t.Errorf(
			"Get() TeamID = %q, want %q",
			got.TeamID,
			wantItem.TeamID,
		)
	}

	if got.AccessToken != wantItem.AccessToken {
		t.Errorf(
			"Get() AccessToken = %q, want %q",
			got.AccessToken,
			wantItem.AccessToken,
		)
	}

	if got.BotUserID != wantItem.BotUserID {
		t.Errorf(
			"Get() BotUserID = %q, want %q",
			got.BotUserID,
			wantItem.BotUserID,
		)
	}

	if client.putItemCalled {
		t.Error("PutItem() was called by Get()")
	}
}
