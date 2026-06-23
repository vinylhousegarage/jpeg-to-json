package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 環境変数を退避
	os.Unsetenv("ALLOWED_ORIGINS")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("AWS_LAMBDA_FUNCTION_NAME")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("S3_BUCKET_NAME")

	t.Run("Success: All required variables set (Local environment)", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000, https://example.com")
		os.Setenv("S3_BUCKET_NAME", "my-test-bucket")
		defer os.Unsetenv("ALLOWED_ORIGINS")
		defer os.Unsetenv("S3_BUCKET_NAME")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		wantOrigins := []string{"http://localhost:3000", "https://example.com"}
		if !reflect.DeepEqual(cfg.AllowedOrigins, wantOrigins) {
			t.Errorf("expected origins %v, got %v", wantOrigins, cfg.AllowedOrigins)
		}
		if cfg.AppEnv != "development" {
      t.Errorf("expected AppEnv 'development', got %s", cfg.AppEnv)
    }
		if cfg.BucketName != "my-test-bucket" {
			t.Errorf("expected bucket name my-test-bucket, got %s", cfg.BucketName)
		}
		if cfg.IsLambda != false {
			t.Errorf("expected IsLambda false, got %v", cfg.IsLambda)
		}
		if cfg.Region != "ap-northeast-1" {
			t.Errorf("expected region ap-northeast-1, got %s", cfg.Region)
		}
	})

	t.Run("Success: Lambda environment", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
		os.Setenv("S3_BUCKET_NAME", "my-test-bucket")
		os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "my-lambda-function")
		defer os.Unsetenv("ALLOWED_ORIGINS")
		defer os.Unsetenv("S3_BUCKET_NAME")
		defer os.Unsetenv("AWS_LAMBDA_FUNCTION_NAME")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.IsLambda != true {
			t.Errorf("expected IsLambda true, got %v", cfg.IsLambda)
		}
	})

	t.Run("Error: Missing ALLOWED_ORIGINS", func(t *testing.T) {
		os.Unsetenv("ALLOWED_ORIGINS")
		os.Setenv("S3_BUCKET_NAME", "test")
		defer os.Unsetenv("S3_BUCKET_NAME")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when ALLOWED_ORIGINS is missing, but got nil")
		}
	})

	t.Run("Error: Missing S3_BUCKET_NAME", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
		os.Unsetenv("S3_BUCKET_NAME")
		defer os.Unsetenv("ALLOWED_ORIGINS")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when S3_BUCKET_NAME is missing, but got nil")
		}
	})
}
