package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AllowedOrigins []string
	AppEnv         string
	BucketName     string
	IsLambda       bool
	Region         string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("backend/.env")
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

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

	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required")
	}

	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	return &Config{
		AllowedOrigins: allowed,
		AppEnv:         env,
		BucketName:     bucket,
		IsLambda:       isLambda,
		Region:         region,
	}, nil
}
