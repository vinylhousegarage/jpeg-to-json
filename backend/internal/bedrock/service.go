package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// BedrockClient のインターフェース定義
type Client interface {
	Invoke(ctx context.Context, input []byte) ([]byte, error)
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{client: client}
}

// ProcessImage は画像データを受け取り、Bedrock経由で構造化データを返す
func (s *Service) ProcessImage(ctx context.Context, rawImage []byte) (map[string]interface{}, error) {
	// 1. Bedrockへ送信 (入力は圧縮済みJPEGバイト)
	response, err := s.client.Invoke(ctx, rawImage)
	if err != nil {
		return nil, fmt.Errorf("bedrock invoke error: %w", err)
	}

	// 2. レスポンスをパース
	return s.parseResponse(response)
}

func (s *Service) parseResponse(response []byte) (map[string]interface{}, error) {
	// Claude/Bedrockの標準的なレスポンス構造
	var resp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bedrock response: %w", err)
	}

	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("no content in bedrock response")
	}

	// JSON文字列部分のクリーニング
	text := resp.Content[0].Text
	text = strings.ReplaceAll(text, "```json", "")
	text = strings.ReplaceAll(text, "```", "")
	text = strings.TrimSpace(text)

	// JSONをパースしてMapへ
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from bedrock output: %w, text: %s", err, text)
	}

	return result, nil
}
