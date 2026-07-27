package storage

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestNewS3Client(t *testing.T) {
	t.Parallel()

	rawClient := &s3.Client{}
	bucket := "my-test-bucket"

	client := NewS3PresignClient(bucket, rawClient)

	if client == nil {
		t.Fatal("expected client to be initialized, got nil")
	}
	if client.bucketName != bucket {
		t.Errorf("expected bucket name %s, got %s", bucket, client.bucketName)
	}
	if client.presignClient == nil {
		t.Error("expected presignClient to be initialized, got nil")
	}
}
