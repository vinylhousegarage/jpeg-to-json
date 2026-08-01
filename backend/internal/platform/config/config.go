package config

import (
	"fmt"
	"os"
)

const (
	appEnvDevelopment = "development"
	appEnvStaging     = "staging"
	appEnvProduction  = "production"
)

type Config struct {
	AppEnv       string
	CookieSecure bool

	IsLambda bool
	Region   string

	BedrockModelID  string
	InputBucketName string
	PromptFileName  string

	SlackClientID     string
	SlackClientSecret string
	SlackRedirectURI  string
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = appEnvDevelopment
	}

	cookieSecure := env == appEnvProduction || env == appEnvStaging

	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		return nil, fmt.Errorf("BEDROCK_MODEL_ID is required")
	}

	inputBucket := os.Getenv("INPUT_BUCKET_NAME")
	if inputBucket == "" {
		return nil, fmt.Errorf("INPUT_BUCKET_NAME is required")
	}

	promptFile := os.Getenv("PROMPT_FILE_NAME")
	if promptFile == "" {
		promptFile = "extractor.txt"
	}

	slackClientID := os.Getenv("SLACK_CLIENT_ID")
	if slackClientID == "" {
		return nil, fmt.Errorf("SLACK_CLIENT_ID is required")
	}

	slackClientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	if slackClientSecret == "" {
		return nil, fmt.Errorf("SLACK_CLIENT_SECRET is required")
	}

	slackRedirectURI := os.Getenv("SLACK_REDIRECT_URI")
	if slackRedirectURI == "" {
		return nil, fmt.Errorf("SLACK_REDIRECT_URI is required")
	}

	return &Config{
		AppEnv:       env,
		CookieSecure: cookieSecure,

		IsLambda: isLambda,
		Region:   region,

		BedrockModelID:  modelID,
		InputBucketName: inputBucket,
		PromptFileName:  promptFile,

		SlackClientID:     slackClientID,
		SlackClientSecret: slackClientSecret,
		SlackRedirectURI:  slackRedirectURI,
	}, nil
}
