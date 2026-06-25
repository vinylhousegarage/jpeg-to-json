package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 環境変数を退避
	_ = os.Unsetenv("ALLOWED_ORIGINS")
	_ = os.Unsetenv("APP_ENV")
	_ = os.Unsetenv("AWS_LAMBDA_FUNCTION_NAME")
	_ = os.Unsetenv("AWS_REGION")
	_ = os.Unsetenv("INPUT_BUCKET_NAME")

	t.Run("Success: All required variables set (Local environment)", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000, https://example.com")
		t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")

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
		if cfg.InputBucketName != "my-test-bucket" {
			t.Errorf("expected bucket name my-test-bucket, got %s", cfg.InputBucketName)
		}
		if cfg.IsLambda != false {
			t.Errorf("expected IsLambda false, got %v", cfg.IsLambda)
		}
		if cfg.Region != "ap-northeast-1" {
			t.Errorf("expected region ap-northeast-1, got %s", cfg.Region)
		}
	})

	t.Run("Success: Lambda environment", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
		t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "my-lambda-function")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.IsLambda != true {
			t.Errorf("expected IsLambda true, got %v", cfg.IsLambda)
		}
	})

	t.Run("Error: Missing ALLOWED_ORIGINS", func(t *testing.T) {
		_ = os.Unsetenv("ALLOWED_ORIGINS")
		t.Setenv("INPUT_BUCKET_NAME", "test")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when ALLOWED_ORIGINS is missing, but got nil")
		}
	})

	t.Run("Error: Missing INPUT_BUCKET_NAME", func(t *testing.T) {
		t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
		_ = os.Unsetenv("INPUT_BUCKET_NAME")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when INPUT_BUCKET_NAME is missing, but got nil")
		}
	})
}
