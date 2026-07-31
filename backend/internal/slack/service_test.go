package slack

import (
	"testing"
)

func TestGenerateState(t *testing.T) {
	t.Parallel()

	state1 := GenerateState()
	state2 := GenerateState()

	if state1 == "" {
		t.Error("GenerateState() returned an empty string")
	}

	if state1 == state2 {
		t.Error("GenerateState() generated duplicate values (not random enough)")
	}
}
