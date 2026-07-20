package bedrock

import "testing"

func TestEncodeBase64(t *testing.T) {
	t.Parallel()

	input := []byte("hello")
	expected := "aGVsbG8="

	result := EncodeBase64(input)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
