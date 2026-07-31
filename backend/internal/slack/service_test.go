package slack

import (
	"net/http"
	"net/url"
	"strings"
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

func TestBuildAuthURL(t *testing.T) {
	t.Parallel()

	clientID := "test-client-id"
	redirectURI := "https://example.com/callback"
	state := "test-state"

	authURLStr := BuildAuthURL(clientID, redirectURI, state)

	parsedURL, err := url.Parse(authURLStr)
	if err != nil {
		t.Fatalf("failed to parse generated auth URL: %v", err)
	}

	if parsedURL.Scheme != "https" || parsedURL.Host != "slack.com" || parsedURL.Path != "/oauth/v2/authorize" {
		t.Errorf("unexpected base URL structure: %s", authURLStr)
	}

	query := parsedURL.Query()

	if query.Get("client_id") != clientID {
		t.Errorf("expected client_id '%s', got '%s'", clientID, query.Get("client_id"))
	}

	if query.Get("scope") != "chat:write,im:write" {
		t.Errorf("expected scope 'chat:write,im:write', got '%s'", query.Get("scope"))
	}

	if query.Get("redirect_uri") != redirectURI {
		t.Errorf("expected redirect_uri '%s', got '%s'", redirectURI, query.Get("redirect_uri"))
	}

	if query.Get("state") != state {
		t.Errorf("expected state '%s', got '%s'", state, query.Get("state"))
	}

	if !strings.Contains(authURLStr, "scope=chat%3Awrite%2Cim%3Awrite") && !strings.Contains(authURLStr, "scope=chat:write,im:write") {
		t.Errorf("scope parameter is malformed: %s", authURLStr)
	}
}
