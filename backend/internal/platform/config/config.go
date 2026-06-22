package config

import (
	"fmt"
	"os"
  "strings"

  "github.com/joho/godotenv"
)

type Config struct {
  AllowedOrigins []string
	BucketName     string
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

	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}

	return &Config{
    AllowedOrigins: allowed,
		BucketName:     bucket,
		Region:         region,
	}, nil
}
