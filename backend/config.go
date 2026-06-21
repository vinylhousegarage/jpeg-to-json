package main

import (
	"fmt"
	"os"

  "github.com/joho/godotenv"
)

type Config struct {
	BucketName string
	Region     string
}

func LoadConfig() (*Config, error) {
  err := godotenv.Load("backend/.env")
    if err != nil {
        fmt.Println("Warning: .env file not found")
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
		BucketName: bucket,
		Region:     region,
	}, nil
}
