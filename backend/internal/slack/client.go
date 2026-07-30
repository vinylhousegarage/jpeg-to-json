package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 構造体を定義
type Client struct {
	shotNumber string
	webhookURL string
	httpClient *http.Client
}

// 構造体を初期化
func NewClient(shotNumber, webhookURL string) *Client {
	return &Client{
		shotNumber: shotNumber,
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ペイロードの構造体を定義
type SlackPayload struct {
	Text string `json:"text"`
}

// 撮影番号・ダウンロードURLをSlackに通知
func (c *Client) SendNotification(ctx context.Context, downloadURL string) error {
	message := fmt.Sprintf("撮影番号: %s\n%s", c.shotNumber, downloadURL)
	payload := SlackPayload{
		Text: message,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack api returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}
