package bedrock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockClient struct {
	sdkClient *bedrockruntime.Client
	modelID   string
}

func NewClient(sdkClient *bedrockruntime.Client, modelID string) *BedrockClient {
	return &BedrockClient{
		sdkClient: sdkClient,
		modelID:   modelID,
	}
}

func (c *BedrockClient) Invoke(ctx context.Context, imgData []byte, prompt string) (*Result, error) {
	return &Result{Data: map[string]string{}}, nil
}
