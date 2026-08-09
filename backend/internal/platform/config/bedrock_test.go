package config

import (
	"strings"
	"testing"
)

func TestLoadBedrockConfig(t *testing.T) {
	t.Run("success with defaults", func(t *testing.T) {
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
			"",
		)

		cfg, err := loadBedrockConfig()
		if err != nil {
			t.Fatalf(
				"loadBedrockConfig() error = %v",
				err,
			)
		}

		if cfg.ModelID != "jp.anthropic.claude-sonnet-4-6" {
			t.Errorf(
				"expected ModelID %q, got %q",
				"jp.anthropic.claude-sonnet-4-6",
				cfg.ModelID,
			)
		}

		if cfg.InputBucketName != "my-test-bucket" {
			t.Errorf(
				"expected InputBucketName %q, got %q",
				"my-test-bucket",
				cfg.InputBucketName,
			)
		}

		if cfg.OutputBucketName != "my-output-bucket" {
			t.Errorf(
				"expected OutputBucketName %q, got %q",
				"my-output-bucket",
				cfg.OutputBucketName,
			)
		}

		if cfg.PromptFileName != defaultPromptFileName {
			t.Errorf(
				"expected PromptFileName %q, got %q",
				defaultPromptFileName,
				cfg.PromptFileName,
			)
		}
	})

	tests := []struct {
		name             string
		modelID          string
		inputBucketName  string
		outputBucketName string
		expectedErr      string
	}{
		{
			name:             "missing model ID",
			modelID:          "",
			inputBucketName:  "my-test-bucket",
			outputBucketName: "my-output-bucket",
			expectedErr:      "BEDROCK_MODEL_ID is required",
		},
		{
			name:             "missing input bucket",
			modelID:          "test-model",
			inputBucketName:  "",
			outputBucketName: "my-output-bucket",
			expectedErr:      "INPUT_BUCKET_NAME is required",
		},
		{
			name:             "missing output bucket",
			modelID:          "test-model",
			inputBucketName:  "my-test-bucket",
			outputBucketName: "",
			expectedErr:      "OUTPUT_BUCKET_NAME is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(
				"BEDROCK_MODEL_ID",
				tt.modelID,
			)
			t.Setenv(
				"INPUT_BUCKET_NAME",
				tt.inputBucketName,
			)
			t.Setenv(
				"OUTPUT_BUCKET_NAME",
				tt.outputBucketName,
			)
			t.Setenv(
				"PROMPT_FILE_NAME",
				"",
			)

			_, err := loadBedrockConfig()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(
				err.Error(),
				tt.expectedErr,
			) {
				t.Errorf(
					"expected error containing %q, got %q",
					tt.expectedErr,
					err.Error(),
				)
			}
		})
	}
}
