package config

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Run("Success: All required variables set (Local environment)", func(t *testing.T) {
		t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")
		t.Setenv("BEDROCK_MODEL_ID", "jp.anthropic.claude-sonnet-4-6")
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_REDIRECT_URI", "https://example.com/callback")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
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
		if cfg.SlackClientID != "test-client-id" {
			t.Errorf("expected SlackClientID 'test-client-id', got %s", cfg.SlackClientID)
		}
	})

	t.Run("Success: Lambda environment", func(t *testing.T) {
		t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")
		t.Setenv("BEDROCK_MODEL_ID", "jp.anthropic.claude-sonnet-4-6")
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "my-lambda-function")
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_REDIRECT_URI", "https://example.com/callback")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.IsLambda != true {
			t.Errorf("expected IsLambda true, got %v", cfg.IsLambda)
		}
	})

	t.Run("Error: Missing INPUT_BUCKET_NAME", func(t *testing.T) {
		t.Setenv("BEDROCK_MODEL_ID", "jp.anthropic.claude-sonnet-4-6")
		t.Setenv("INPUT_BUCKET_NAME", "")
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_REDIRECT_URI", "https://example.com/callback")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when INPUT_BUCKET_NAME is missing, but got nil")
		}
	})

	t.Run("Error: Missing BEDROCK_MODEL_ID", func(t *testing.T) {
		t.Setenv("INPUT_BUCKET_NAME", "test")
		t.Setenv("BEDROCK_MODEL_ID", "")
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_REDIRECT_URI", "https://example.com/callback")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when BEDROCK_MODEL_ID is missing, but got nil")
		}
	})

	t.Run("Error: Missing SLACK_CLIENT_ID", func(t *testing.T) {
		t.Setenv("INPUT_BUCKET_NAME", "test")
		t.Setenv("BEDROCK_MODEL_ID", "test")
		t.Setenv("SLACK_CLIENT_ID", "")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv("SLACK_REDIRECT_URI", "https://example.com/callback")

		_, err := LoadConfig()
		if err == nil {
			t.Error("expected error when SLACK_CLIENT_ID is missing, but got nil")
		}
	})
}
