package bedrock

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockRuntimeClient interface {
	InvokeModel(
		ctx context.Context,
		params *bedrockruntime.InvokeModelInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.InvokeModelOutput, error)
}

type BedrockClient struct {
	sdkClient BedrockRuntimeClient
	modelID   string
	prompt    string
	logger    *zap.Logger
}

func NewClient(
	ctx context.Context,
	modelID string,
	prompt string,
	logger *zap.Logger,
) (*BedrockClient, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRetryMaxAttempts(5),
		config.WithRetryMode(aws.RetryModeStandard),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &BedrockClient{
		sdkClient: bedrockruntime.NewFromConfig(cfg),
		modelID:   modelID,
		prompt:    prompt,
		logger:    logger,
	}, nil
}

func (c *BedrockClient) Analyze(ctx context.Context, imgData []byte) (*ResponseBody, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imgData)

	payload := RequestBody{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        2000,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{Type: "image", Source: &ImageSource{Type: "base64", MediaType: "image/jpeg", Data: imageBase64}},
					{Type: "text", Text: c.prompt},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     &c.modelID,
		ContentType: aws.String("application/json"),
		Body:        body,
	}

	output, err := c.sdkClient.InvokeModel(ctx, input)
	if err != nil {
		return nil, err
	}

	// レスポンス（output.Body）を ResponseBody 構造体でパース
	var resp ResponseBody
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bedrock response: %w", err)
	}

	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("no content in bedrock response")
	}

	// ログ出力
	c.logger.Info("bedrock inference success",
		zap.Int("input_tokens", resp.Usage.InputTokens),
		zap.Int("output_tokens", resp.Usage.OutputTokens),
	)

	return &resp, nil
}
