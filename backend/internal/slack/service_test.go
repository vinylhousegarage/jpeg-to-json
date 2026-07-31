package slack

import (
	"net/http"
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

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	testState := "test-random-state-string"
	cookie := BuildStateCookie(testState)

	if cookie.Name != "oauth_state" {
		t.Errorf("expected cookie name 'oauth_state', got '%s'", cookie.Name)
	}

	if cookie.Value != testState {
		t.Errorf("expected cookie value '%s', got '%s'", testState, cookie.Value)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("expected Secure to be true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite LaxMode, got %v", cookie.SameSite)
	}

	if cookie.Path != "/" {
		t.Errorf("expected Path '/', got '%s'", cookie.Path)
	}
}
