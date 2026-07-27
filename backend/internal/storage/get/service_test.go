package get

import (
	"context"
	"errors"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// モック用の構造体
type mockPresigner struct {
	mockResult *v4.PresignedHTTPRequest
	mockErr    error
}

func (m *mockPresigner) PresignGetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	return m.mockResult, m.mockErr
}

// 正常系
func TestService_GeneratePresignURL_Success(t *testing.T) {
	t.Parallel()

	expectedURL := "https://example.com/presigned-url"
	mock := &mockPresigner{
		mockResult: &v4.PresignedHTTPRequest{
			URL: expectedURL,
		},
		mockErr: nil,
	}

	svc := NewService(mock)

	url, expiresAt, err := svc.GeneratePresignURL(context.Background(), "test.jpg")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != expectedURL {
		t.Errorf("expected URL %s, got %s", expectedURL, url)
	}

	if time.Until(expiresAt) > 16*time.Minute || time.Until(expiresAt) < 14*time.Minute {
		t.Errorf("unexpected expiration time: %v", expiresAt)
	}
}

// 異常系
func TestService_GeneratePresignURL_Error(t *testing.T) {
	t.Parallel()

	mock := &mockPresigner{
		mockResult: nil,
		mockErr:    errors.New("aws error"),
	}

	svc := NewService(mock)

	_, _, err := svc.GeneratePresignURL(context.Background(), "test.jpg")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
