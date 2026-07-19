package bedrock

import (
	"strings"
	"testing"
)

func TestLoadPrompt(t *testing.T) {
	t.Parallel()

	content, err := LoadPrompt()
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	if content == "" {
		t.Error("Loaded prompt is empty")
	}

	if !strings.Contains(content, "<system_role>") {
		t.Error("Prompt does not contain expected system_role tag")
	}
}

func TestEncodeBase64(t *testing.T) {
	t.Parallel()

	input := []byte("hello")
	expected := "aGVsbG8="

	result := EncodeBase64(input)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
