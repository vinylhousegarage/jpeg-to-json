package presign

import (
	"testing"
)

func TestNewPresignClient(t *testing.T) {
	t.Parallel()

	mock := &mockPresigner{}
	bucket := "my-test-bucket"

	client := NewPresignClient(mock, bucket)

	if client == nil {
		t.Fatal("expected client to be initialized, got nil")
	}
	if client.s3Presigner != mock {
		t.Error("expected s3Presigner to be the provided mock")
	}
	if client.bucketName != bucket {
		t.Errorf("expected bucket name %s, got %s", bucket, client.bucketName)
	}
}
