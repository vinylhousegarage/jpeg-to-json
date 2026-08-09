package config

import (
	"fmt"
	"os"
)

const (
	appEnvDevelopment = "development"
	appEnvStaging     = "staging"
	appEnvProduction  = "production"

	defaultAWSRegion      = "ap-northeast-1"
	defaultPromptFileName = "extractor.txt"
)

type Config struct {
	App     AppConfig
	AWS     AWSConfig
	Bedrock BedrockConfig
	Slack   SlackConfig
}

type AppConfig struct {
	Env          string
	CookieSecure bool
}

type AWSConfig struct {
	IsLambda bool
	Region   string
}

type BedrockConfig struct {
	ModelID         string
	InputBucketName string
	PromptFileName  string
}

type SlackConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func Load() (*Config, error) {
	appConfig := loadAppConfig()
	awsConfig := loadAWSConfig()

	bedrockConfig, err := loadBedrockConfig()
	if err != nil {
		return nil, err
	}

	slackConfig, err := loadSlackConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		App:     appConfig,
		AWS:     awsConfig,
		Bedrock: bedrockConfig,
		Slack:   slackConfig,
	}, nil
}

func loadAppConfig() AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = appEnvDevelopment
	}

	return AppConfig{
		Env: env,
		CookieSecure: env == appEnvProduction ||
			env == appEnvStaging,
	}
}

func loadAWSConfig() AWSConfig {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultAWSRegion
	}

	return AWSConfig{
		IsLambda: os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "",
		Region:   region,
	}
}

func loadBedrockConfig() (BedrockConfig, error) {
	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		return BedrockConfig{}, fmt.Errorf(
			"BEDROCK_MODEL_ID is required",
		)
	}

	inputBucketName := os.Getenv("INPUT_BUCKET_NAME")
	if inputBucketName == "" {
		return BedrockConfig{}, fmt.Errorf(
			"INPUT_BUCKET_NAME is required",
		)
	}

	promptFileName := os.Getenv("PROMPT_FILE_NAME")
	if promptFileName == "" {
		promptFileName = defaultPromptFileName
	}

	return BedrockConfig{
		ModelID:         modelID,
		InputBucketName: inputBucketName,
		PromptFileName:  promptFileName,
	}, nil
}

func loadSlackConfig() (SlackConfig, error) {
	clientID := os.Getenv("SLACK_CLIENT_ID")
	if clientID == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_CLIENT_ID is required",
		)
	}

	clientSecret := os.Getenv("SLACK_CLIENT_SECRET")
	if clientSecret == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_CLIENT_SECRET is required",
		)
	}

	redirectURI := os.Getenv("SLACK_REDIRECT_URI")
	if redirectURI == "" {
		return SlackConfig{}, fmt.Errorf(
			"SLACK_REDIRECT_URI is required",
		)
	}

	return SlackConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}, nil
}
