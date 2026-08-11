package notifier

import (
	"context"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

type stubTokenStore struct {
	token *oauth.Token
	err   error

	called bool
	ctx    context.Context
}

func (s *stubTokenStore) Get(
	ctx context.Context,
) (*oauth.Token, error) {
	s.called = true
	s.ctx = ctx

	return s.token, s.err
}

type stubMessageClient struct {
	err error

	called      bool
	ctx         context.Context
	accessToken string
	channelID   string
	text        string
}

func (s *stubMessageClient) PostMessage(
	ctx context.Context,
	accessToken string,
	channelID string,
	text string,
) error {
	s.called = true
	s.ctx = ctx
	s.accessToken = accessToken
	s.channelID = channelID
	s.text = text

	return s.err
}
