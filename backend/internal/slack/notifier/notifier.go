package notifier

import (
	"context"
	"fmt"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

type tokenStore interface {
	Get(
		ctx context.Context,
		teamID string,
	) (*oauth.Token, error)
}

type messageClient interface {
	PostMessage(
		ctx context.Context,
		accessToken string,
		channelID string,
		text string,
	) error
}

// Notifier sends Slack messages using stored OAuth tokens.
type Notifier struct {
	tokenStore tokenStore
	client     messageClient
}

func NewNotifier(
	tokenStore tokenStore,
	client messageClient,
) *Notifier {
	return &Notifier{
		tokenStore: tokenStore,
		client:     client,
	}
}

func (n *Notifier) Notify(
	ctx context.Context,
	teamID string,
	text string,
) error {
	if teamID == "" {
		return fmt.Errorf("notify slack: team ID is empty")
	}

	if text == "" {
		return fmt.Errorf("notify slack: text is empty")
	}

	token, err := n.tokenStore.Get(ctx, teamID)
	if err != nil {
		return fmt.Errorf("get slack token: %w", err)
	}

	if token == nil {
		return fmt.Errorf("get slack token: token is nil")
	}

	if token.AccessToken == "" {
		return fmt.Errorf("get slack token: access token is empty")
	}

	if token.ChannelID == "" {
		return fmt.Errorf("get slack token: channel ID is empty")
	}

	if err := n.client.PostMessage(
		ctx,
		token.AccessToken,
		token.ChannelID,
		text,
	); err != nil {
		return fmt.Errorf("post slack message: %w", err)
	}

	return nil
}
