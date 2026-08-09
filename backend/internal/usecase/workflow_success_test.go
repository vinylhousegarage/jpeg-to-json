package usecase

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func TestWorkflow_Execute_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := zap.NewNop()

	getOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(
			bytes.NewReader([]byte("fake-image-bytes")),
		),
	}

	mockGetter := &mockS3Getter{
		getOutput: getOutput,
	}

	mockPutter := &mockS3Putter{}

	expectedDownloadURL := "https://s3.amazonaws.com/mock-presigned-url"
	mockPresigner := &mockS3Presigner{
		presignURL: expectedDownloadURL,
	}

	mockBedrock := &mockBedrockService{
		resultMap: map[string]any{
			"key": "value",
		},
	}

	mockSlack := &mockSlackClient{}

	workflow := NewWorkflow(
		mockGetter,
		mockPutter,
		mockPresigner,
		mockBedrock,
		mockSlack,
		"output-bucket",
		logger,
	)

	err := workflow.Execute(
		ctx,
		"input-bucket",
		"SHOT-001.jpg",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mockSlack.called {
		t.Fatal("expected slack notification to be called")
	}

	if mockSlack.shot != "SHOT-001" {
		t.Errorf(
			"expected shot number %q, got %q",
			"SHOT-001",
			mockSlack.shot,
		)
	}

	if mockSlack.msg != expectedDownloadURL {
		t.Errorf(
			"expected download URL %q, got %q",
			expectedDownloadURL,
			mockSlack.msg,
		)
	}
}
