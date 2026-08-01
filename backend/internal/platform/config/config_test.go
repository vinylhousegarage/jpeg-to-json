package config

import (
	"strings"
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()

	t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")
	t.Setenv("BEDROCK_MODEL_ID", "jp.anthropic.claude-sonnet-4-6")
	t.Setenv("SLACK_CLIENT_ID", "test-client-id")
	t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
	t.Setenv("SLACK_REDIRECT_URI", "https://example.com/oauth/slack/callback")
}

func TestLoadConfig(t *testing.T) {
	t.Run("Success: development environment with default values", func(t *testing.T) {
		setRequiredEnv(t)

		t.Setenv("APP_ENV", "")
		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")
		t.Setenv("AWS_REGION", "")
		t.Setenv("PROMPT_FILE_NAME", "")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.AppEnv != appEnvDevelopment {
			t.Errorf(
				"expected AppEnv %q, got %q",
				appEnvDevelopment,
				cfg.AppEnv,
			)
		}

		if cfg.CookieSecure {
			t.Error("expected CookieSecure false in development")
		}

		if cfg.IsLambda {
			t.Error("expected IsLambda false")
		}

		if cfg.Region != "ap-northeast-1" {
			t.Errorf(
				"expected Region %q, got %q",
				"ap-northeast-1",
				cfg.Region,
			)
		}

		if cfg.BedrockModelID != "jp.anthropic.claude-sonnet-4-6" {
			t.Errorf(
				"expected BedrockModelID %q, got %q",
				"jp.anthropic.claude-sonnet-4-6",
				cfg.BedrockModelID,
			)
		}

		if cfg.InputBucketName != "my-test-bucket" {
			t.Errorf(
				"expected InputBucketName %q, got %q",
				"my-test-bucket",
				cfg.InputBucketName,
			)
		}

		if cfg.PromptFileName != "extractor.txt" {
			t.Errorf(
				"expected PromptFileName %q, got %q",
				"extractor.txt",
				cfg.PromptFileName,
			)
		}

		if cfg.SlackClientID != "test-client-id" {
			t.Errorf(
				"expected SlackClientID %q, got %q",
				"test-client-id",
				cfg.SlackClientID,
			)
		}

		if cfg.SlackClientSecret != "test-client-secret" {
			t.Errorf(
				"expected SlackClientSecret %q, got %q",
				"test-client-secret",
				cfg.SlackClientSecret,
			)
		}

		expectedRedirectURI := "https://example.com/oauth/slack/callback"
		if cfg.SlackRedirectURI != expectedRedirectURI {
			t.Errorf(
				"expected SlackRedirectURI %q, got %q",
				expectedRedirectURI,
				cfg.SlackRedirectURI,
			)
		}
	})

	t.Run("Success: production environment", func(t *testing.T) {
		setRequiredEnv(t)

		t.Setenv("APP_ENV", appEnvProduction)
		t.Setenv("AWS_REGION", "us-east-1")
		t.Setenv("PROMPT_FILE_NAME", "custom-prompt.txt")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg.AppEnv != appEnvProduction {
			t.Errorf(
				"expected AppEnv %q, got %q",
				appEnvProduction,
				cfg.AppEnv,
			)
		}

		if !cfg.CookieSecure {
			t.Error("expected CookieSecure true in production")
		}

		if cfg.Region != "us-east-1" {
			t.Errorf(
				"expected Region %q, got %q",
				"us-east-1",
				cfg.Region,
			)
		}

		if cfg.PromptFileName != "custom-prompt.txt" {
			t.Errorf(
				"expected PromptFileName %q, got %q",
				"custom-prompt.txt",
				cfg.PromptFileName,
			)
		}
	})

	t.Run("Success: staging environment", func(t *testing.T) {
		setRequiredEnv(t)

		t.Setenv("APP_ENV", appEnvStaging)

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !cfg.CookieSecure {
			t.Error("expected CookieSecure true in staging")
		}
	})

	t.Run("Success: Lambda environment", func(t *testing.T) {
		setRequiredEnv(t)

		t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "my-lambda-function")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !cfg.IsLambda {
			t.Error("expected IsLambda true")
		}
	})

	tests := []struct {
		name        string
		missingKey  string
		expectedErr string
	}{
		{
			name:        "Missing BEDROCK_MODEL_ID",
			missingKey:  "BEDROCK_MODEL_ID",
			expectedErr: "BEDROCK_MODEL_ID is required",
		},
		{
			name:        "Missing INPUT_BUCKET_NAME",
			missingKey:  "INPUT_BUCKET_NAME",
			expectedErr: "INPUT_BUCKET_NAME is required",
		},
		{
			name:        "Missing SLACK_CLIENT_ID",
			missingKey:  "SLACK_CLIENT_ID",
			expectedErr: "SLACK_CLIENT_ID is required",
		},
		{
			name:        "Missing SLACK_CLIENT_SECRET",
			missingKey:  "SLACK_CLIENT_SECRET",
			expectedErr: "SLACK_CLIENT_SECRET is required",
		},
		{
			name:        "Missing SLACK_REDIRECT_URI",
			missingKey:  "SLACK_REDIRECT_URI",
			expectedErr: "SLACK_REDIRECT_URI is required",
		},
	}

	for _, tt := range tests {
		t.Run("Error: "+tt.name, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv(tt.missingKey, "")

			_, err := LoadConfig()
			if err == nil {
				t.Fatalf(
					"expected error when %s is missing, got nil",
					tt.missingKey,
				)
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf(
					"expected error containing %q, got %q",
					tt.expectedErr,
					err.Error(),
				)
			}
		})
	}
}
