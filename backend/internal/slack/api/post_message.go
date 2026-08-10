package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type postMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type postMessageResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func (c *Client) PostMessage(
	ctx context.Context,
	accessToken string,
	channelID string,
	text string,
) error {
	if accessToken == "" {
		return fmt.Errorf(
			"post slack message: access token is empty",
		)
	}

	if channelID == "" {
		return fmt.Errorf(
			"post slack message: channel ID is empty",
		)
	}

	if text == "" {
		return fmt.Errorf(
			"post slack message: text is empty",
		)
	}

	body, err := json.Marshal(
		postMessageRequest{
			Channel: channelID,
			Text:    text,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"marshal slack post message request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.postMessageURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf(
			"create slack post message request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+accessToken,
	)
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"send slack post message request: %w",
			err,
		)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"slack post message request returned status %d",
			res.StatusCode,
		)
	}

	var payload postMessageResponse
	if err := json.NewDecoder(
		res.Body,
	).Decode(&payload); err != nil {
		return fmt.Errorf(
			"decode slack post message response: %w",
			err,
		)
	}

	if !payload.OK {
		return fmt.Errorf(
			"slack post message failed: %s",
			payload.Error,
		)
	}

	return nil
}
