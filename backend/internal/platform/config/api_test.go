package config

import "testing"

func TestLoadAPI(t *testing.T) {
	t.Setenv(
		"APP_ENV",
		appEnvLocal,
	)
	t.Setenv(
		"AWS_LAMBDA_FUNCTION_NAME",
		"",
	)
	t.Setenv(
		"AWS_REGION",
		defaultAWSRegion,
	)

	t.Setenv(
		"INPUT_BUCKET_NAME",
		"my-test-bucket",
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

	cfg, err := LoadAPI()
	if err != nil {
		t.Fatalf(
			"LoadAPI() error = %v",
			err,
		)
	}

	if cfg.App.Env != appEnvLocal {
		t.Errorf(
			"expected App.Env %q, got %q",
			appEnvLocal,
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

	if cfg.Storage.InputBucketName != "my-test-bucket" {
		t.Errorf(
			"expected Storage.InputBucketName %q, got %q",
			"my-test-bucket",
			cfg.Storage.InputBucketName,
		)
	}

	if cfg.Storage.OutputBucketName != "" {
		t.Errorf(
			"expected Storage.OutputBucketName to be empty, got %q",
			cfg.Storage.OutputBucketName,
		)
	}

	if cfg.Slack.ClientID != "test-client-id" {
		t.Errorf(
			"expected Slack.ClientID %q, got %q",
			"test-client-id",
			cfg.Slack.ClientID,
		)
	}

	if cfg.Slack.ClientSecret != "test-client-secret" {
		t.Errorf(
			"expected Slack.ClientSecret %q, got %q",
			"test-client-secret",
			cfg.Slack.ClientSecret,
		)
	}

	if cfg.Slack.RedirectURI != "https://example.com/oauth/slack/callback" {
		t.Errorf(
			"expected Slack.RedirectURI %q, got %q",
			"https://example.com/oauth/slack/callback",
			cfg.Slack.RedirectURI,
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
