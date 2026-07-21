package bedrock

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

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
}

func NewClient(
	ctx context.Context,
	modelID string,
) (*BedrockClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &BedrockClient{
		sdkClient: bedrockruntime.NewFromConfig(cfg),
		modelID:   modelID,
	}, nil
}

func (c *BedrockClient) Invoke(ctx context.Context, imgData []byte, prompt string) (*Result, error) {
	imageBase64 := base64.StdEncoding.EncodeToString(imgData)

	payload := RequestBody{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        2000,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{Type: "image", Source: &ImageSource{Type: "base64", MediaType: "image/jpeg", Data: imageBase64}},
					{Type: "text", Text: prompt},
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

	// 1. レスポンス（output.Body）を ResponseBody 構造体でパース
	var resp ResponseBody
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bedrock response: %w", err)
	}

	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("no content in bedrock response")
	}

	// 2. Claude が返した JSON 文字列（resp.Content[0].Text）をパース
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Content[0].Text), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON string: %w, text: %s", err, resp.Content[0].Text)
	}

	// 3. Result 型に入れて返す
	return &Result{Data: data}, nil
}
