package bedrock

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type MockBedrockRuntimeClient struct {
	InvokeModelFunc func(
		ctx context.Context,
		params *bedrockruntime.InvokeModelInput,
		optFns ...func(*bedrockruntime.Options),
	) (*bedrockruntime.InvokeModelOutput, error)
}

func (m *MockBedrockRuntimeClient) InvokeModel(
	ctx context.Context,
	params *bedrockruntime.InvokeModelInput,
	optFns ...func(*bedrockruntime.Options),
) (*bedrockruntime.InvokeModelOutput, error) {
	return m.InvokeModelFunc(ctx, params, optFns...)
}

func TestAnalyze(t *testing.T) {
	t.Parallel()

	mockClient := &MockBedrockRuntimeClient{
		InvokeModelFunc: func(
			ctx context.Context,
			params *bedrockruntime.InvokeModelInput,
			optFns ...func(*bedrockruntime.Options),
		) (*bedrockruntime.InvokeModelOutput, error) {
			responseJSON := `{
				"content": [{"text": "{\"key\": \"value\"}"}],
				"usage": {
					"input_tokens": 100,
					"output_tokens": 50
				}
			}`

			return &bedrockruntime.InvokeModelOutput{
				Body: []byte(responseJSON),
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

	expectedContent := "{\"key\": \"value\"}"
	if len(result.Content) == 0 {
		t.Fatalf("expected content to have at least one element")
	}
	if result.Content[0].Text != expectedContent {
		t.Errorf("expected content text to be %s, got %s", expectedContent, result.Content[0].Text)
	}
}
