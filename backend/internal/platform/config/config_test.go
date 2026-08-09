package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("APP_ENV", appEnvDevelopment)
	t.Setenv("AWS_LAMBDA_FUNCTION_NAME", "")
	t.Setenv("AWS_REGION", defaultAWSRegion)
	t.Setenv("BEDROCK_MODEL_ID", "jp.anthropic.claude-sonnet-4-6")
	t.Setenv("INPUT_BUCKET_NAME", "my-test-bucket")
	t.Setenv("PROMPT_FILE_NAME", defaultPromptFileName)
	t.Setenv("SLACK_CLIENT_ID", "test-client-id")
	t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
	t.Setenv(
		"SLACK_REDIRECT_URI",
		"https://example.com/oauth/slack/callback",
	)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
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

	if cfg.Slack.ClientID != "test-client-id" {
		t.Errorf(
			"expected Slack.ClientID %q, got %q",
			"test-client-id",
			cfg.Slack.ClientID,
		)
	}
}
