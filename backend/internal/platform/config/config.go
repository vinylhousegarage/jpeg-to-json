package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AllowedOrigins  []string
	AppEnv          string
	BedrockModelID  string
	InputBucketName string
	IsLambda        bool
	PromptFileName  string
	Region          string
}

func LoadConfig() (*Config, error) {

	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return nil, fmt.Errorf("ALLOWED_ORIGINS is required")
	}

	var allowed []string
	for _, s := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			allowed = append(allowed, trimmed)
		}
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		return nil, fmt.Errorf("BEDROCK_MODEL_ID is required")
	}

	inputBucket := os.Getenv("INPUT_BUCKET_NAME")
	if inputBucket == "" {
		return nil, fmt.Errorf("INPUT_BUCKET_NAME is required")
	}

	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""

	promptFile := os.Getenv("PROMPT_FILE_NAME")
	if promptFile == "" {
		promptFile = "extractor.txt"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	return &Config{
		AllowedOrigins:  allowed,
		AppEnv:          env,
		BedrockModelID:  modelID,
		InputBucketName: inputBucket,
		IsLambda:        isLambda,
		PromptFileName:  promptFile,
		Region:          region,
	}, nil
}
