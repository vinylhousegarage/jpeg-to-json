package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

func TestWorkflow_Execute_GetObjectError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getErr: errors.New("s3 get error"),
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to get object from s3") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_ReadImageError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: &errorReadCloser{
					err: errors.New("read error"),
				},
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read image body") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_BedrockError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{
			processErr: errors.New("bedrock error"),
		},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to process image with bedrock",
	) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_MarshalError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{},
		&mockBedrockService{
			resultMap: map[string]any{
				"invalid": make(chan int),
			},
		},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to marshal json") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_PutObjectError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{
			putErr: errors.New("s3 put error"),
		},
		&mockS3Presigner{},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to put json to output s3",
	) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_PresignError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{
			presignErr: errors.New("presign error"),
		},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		&mockSlackClient{},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to generate presigned url",
	) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflow_Execute_SlackError(t *testing.T) {
	t.Parallel()

	workflow := NewWorkflow(
		&mockS3Getter{
			getOutput: &s3.GetObjectOutput{
				Body: io.NopCloser(
					bytes.NewReader([]byte("fake-image-bytes")),
				),
			},
		},
		&mockS3Putter{},
		&mockS3Presigner{
			presignURL: "https://example.com/download",
		},
		&mockBedrockService{
			resultMap: map[string]any{
				"key": "value",
			},
		},
		&mockSlackClient{
			sendErr: errors.New("slack error"),
		},
		"output-bucket",
		zap.NewNop(),
	)

	err := workflow.Execute(
		context.Background(),
		"input-bucket",
		"SHOT-001.jpg",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(
		err.Error(),
		"failed to send slack notification",
	) {
		t.Errorf("unexpected error: %v", err)
	}
}
