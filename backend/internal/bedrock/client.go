package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient struct {
	sdkClient *bedrockruntime.Client
	modelID   string
}

func NewClient(ctx context.Context, modelID string) (*BedrockClient, error) {
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
	return &Result{Data: map[string]string{}}, nil
}
