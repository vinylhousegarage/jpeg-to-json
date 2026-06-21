package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
  // 環境変数を退避
    os.Unsetenv("S3_BUCKET_NAME")
    os.Unsetenv("AWS_REGION")

	// 正常形
	os.Setenv("S3_BUCKET_NAME", "my-test-bucket")
	defer os.Unsetenv("S3_BUCKET_NAME")

	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if cfg.BucketName != "my-test-bucket" {
		t.Errorf("expected bucket name my-test-bucket, got %s", cfg.BucketName)
	}

	// 異常系
	os.Unsetenv("S3_BUCKET_NAME")
	_, err = LoadConfig()
	if err == nil {
		t.Error("expected error when S3_BUCKET_NAME is missing, but got nil")
	}
}
