package bedrock

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type MockBedrockRuntimeClient struct {
	ConverseFunc func(
		ctx context.Context,
		params *bedrockruntime.ConverseInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.ConverseOutput, error)
}

func (m *MockBedrockRuntimeClient) Converse(
	ctx context.Context,
	params *bedrockruntime.ConverseInput,
	optFns ...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseOutput, error) {
	return m.ConverseFunc(ctx, params, optFns...)
}

func TestAnalyze(t *testing.T) {
	t.Parallel()

	expectedText := `{"key": "value"}`

	mockClient := &MockBedrockRuntimeClient{
		ConverseFunc: func(
			ctx context.Context,
			params *bedrockruntime.ConverseInput,
			optFns ...func(*bedrockruntime.Options),
		) (*bedrockruntime.ConverseOutput, error) {
			return &bedrockruntime.ConverseOutput{
				Output: &types.ConverseOutputMemberMessage{
					Value: types.Message{
						Role: types.ConversationRoleAssistant,
						Content: []types.ContentBlock{
							&types.ContentBlockMemberText{
								Value: expectedText,
							},
						},
					},
				},
				Usage: &types.TokenUsage{
					InputTokens:  aws.Int32(100),
					OutputTokens: aws.Int32(50),
				},
			}, nil
		},
	}

	client := &BedrockClient{
		sdkClient: mockClient,
		modelID:   "test-model",
		prompt:    "extract test",
		logger:    zap.NewNop(),
	}

	result, err := client.Analyze(context.Background(), []byte("fake-image-data"))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expectedText {
		t.Errorf("expected content text to be %s, got %s", expectedText, result)
	}
}
