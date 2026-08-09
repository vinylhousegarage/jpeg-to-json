package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", appEnvDevelopment)
	t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")
	t.Setenv("AWS_REGION", defaultAWSRegion)

	t.Setenv(
		"BEDROCK_MODEL_ID",
		"jp.anthropic.claude-sonnet-4-6",
	)
	t.Setenv(
		"INPUT_BUCKET_NAME",
		"my-test-bucket",
	)
	t.Setenv(
		"OUTPUT_BUCKET_NAME",
		"my-output-bucket",
	)
	t.Setenv(
		"PROMPT_FILE_NAME",
		defaultPromptFileName,
	)

	t.Setenv(
		"SLACK_CLIENT_ID",
		"test-client-id",
	)
	t.Setenv(
		"SLACK_CLIENT_SECRET",
		"test-client-secret",
	)
	t.Setenv(
		"SLACK_REDIRECT_URI",
		"https://example.com/oauth/slack/callback",
	)
	t.Setenv(
		"SLACK_TOKEN_TABLE_NAME",
		"slack-tokens",
	)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.App.Env != appEnvDevelopment {
		t.Errorf(
			"expected App.Env %q, got %q",
			appEnvDevelopment,
			cfg.App.Env,
		)
	}

	if cfg.AWS.Region != defaultAWSRegion {
		t.Errorf(
			"expected AWS.Region %q, got %q",
			defaultAWSRegion,
			cfg.AWS.Region,
		)
	}

	if cfg.Bedrock.ModelID != "jp.anthropic.claude-sonnet-4-6" {
		t.Errorf(
			"expected Bedrock.ModelID %q, got %q",
			"jp.anthropic.claude-sonnet-4-6",
			cfg.Bedrock.ModelID,
		)
	}

	if cfg.Bedrock.OutputBucketName != "my-output-bucket" {
		t.Errorf(
			"expected Bedrock.OutputBucketName %q, got %q",
			"my-output-bucket",
			cfg.Bedrock.OutputBucketName,
		)
	}

	if cfg.Slack.ClientID != "test-client-id" {
		t.Errorf(
			"expected Slack.ClientID %q, got %q",
			"test-client-id",
			cfg.Slack.ClientID,
		)
	}

	if cfg.Slack.TokenTableName != "slack-tokens" {
		t.Errorf(
			"expected Slack.TokenTableName %q, got %q",
			"slack-tokens",
			cfg.Slack.TokenTableName,
		)
	}
}
