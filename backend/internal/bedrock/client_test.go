package bedrock

import (
	"context"
	"testing"

	
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type MockBedrockRuntimeClient struct {
	// InvokeModel を保存してテストで検証可能にする
	InvokeModelFunc func(
		ctx	context.Context,
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

func TestInvoke(t *testing.T) {
	mockClient := &MockBedrockRuntimeClient{
		InvokeModelFunc: func(
			ctx context.Context,
			params *bedrockruntime.InvokeModelInput,
			optFns ...func(*bedrockruntime.Options),
		) (*bedrockruntime.InvokeModelOutput, error) {
			// Bedrockのレスポンスを模倣
			return &bedrockruntime.InvokeModelOutput{
				Body: []byte(`{"content": [{"text": "{\"key\": \"value\"}"}]}`),
			}, nil
		},
	}

	client := &BedrockClient{sdkClient: mockClient, modelID: "test-model"}
	
	// 実行
	result, err := client.Invoke(context.Background(), []byte("fake-image-data"), "extract test")
	
	// 検証
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Result型の中身（Data map[string]string）を検証
	expectedVal := "value"
	if result.Data["key"] != expectedVal {
		t.Errorf("expected key to be %s, got %s", expectedVal, result.Data["key"])
	}
}
