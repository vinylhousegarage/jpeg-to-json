package oauth

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestGenerateState(t *testing.T) {
	t.Parallel()

	state1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}

	state2, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() returned error: %v", err)
	}

	if state1 == "" {
		t.Error("GenerateState() returned an empty string")
	}

	if state1 == state2 {
		t.Error("GenerateState() generated duplicate values")
	}
}

func TestBuildStateCookie(t *testing.T) {
	t.Parallel()

	testState := "test-random-state-string"
	before := time.Now()

	cookie := BuildStateCookie(testState, true)

	if cookie.Name != oauthStateCookieName {
		t.Errorf(
			"expected cookie name %q, got %q",
			oauthStateCookieName,
			cookie.Name,
		)
	}

	if cookie.Value != testState {
		t.Errorf(
			"expected cookie value %q, got %q",
			testState,
			cookie.Value,
		)
	}

	if cookie.Path != oauthStateCookiePath {
		t.Errorf(
			"expected Path %q, got %q",
			oauthStateCookiePath,
			cookie.Path,
		)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("expected Secure to be true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"expected SameSite LaxMode, got %v",
			cookie.SameSite,
		)
	}

	expectedMaxAge := int(OAuthStateTTL.Seconds())
	if cookie.MaxAge != expectedMaxAge {
		t.Errorf(
			"expected MaxAge %d, got %d",
			expectedMaxAge,
			cookie.MaxAge,
		)
	}

	minExpires := before.Add(OAuthStateTTL)
	maxExpires := time.Now().Add(OAuthStateTTL)

	if cookie.Expires.Before(minExpires) || cookie.Expires.After(maxExpires) {
		t.Errorf(
			"expected Expires between %v and %v, got %v",
			minExpires,
			maxExpires,
			cookie.Expires,
		)
	}
}

func TestBuildStateCookie_NotSecure(t *testing.T) {
	t.Parallel()

	cookie := BuildStateCookie("test-state", false)

	if cookie.Secure {
		t.Error("expected Secure to be false")
	}
}

func TestBuildDeleteStateCookie(t *testing.T) {
	t.Parallel()

	cookie := BuildDeleteStateCookie(true)

	if cookie.Name != oauthStateCookieName {
		t.Errorf(
			"expected cookie name %q, got %q",
			oauthStateCookieName,
			cookie.Name,
		)
	}

	if cookie.Value != "" {
		t.Errorf("expected empty cookie value, got %q", cookie.Value)
	}

	if cookie.Path != oauthStateCookiePath {
		t.Errorf(
			"expected Path %q, got %q",
			oauthStateCookiePath,
			cookie.Path,
		)
	}

	if !cookie.HttpOnly {
		t.Error("expected HttpOnly to be true")
	}

	if !cookie.Secure {
		t.Error("expected Secure to be true")
	}

	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"expected SameSite LaxMode, got %v",
			cookie.SameSite,
		)
	}

	if cookie.MaxAge != -1 {
		t.Errorf("expected MaxAge -1, got %d", cookie.MaxAge)
	}

	if !cookie.Expires.Equal(time.Unix(0, 0)) {
		t.Errorf(
			"expected Unix epoch expiration, got %v",
			cookie.Expires,
		)
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

	if parsedURL.Scheme != "https" {
		t.Errorf("expected scheme https, got %q", parsedURL.Scheme)
	}

	if parsedURL.Host != "slack.com" {
		t.Errorf("expected host slack.com, got %q", parsedURL.Host)
	}

	if parsedURL.Path != "/oauth/v2/authorize" {
		t.Errorf(
			"expected path %q, got %q",
			"/oauth/v2/authorize",
			parsedURL.Path,
		)
	}

	query := parsedURL.Query()

	if query.Get("client_id") != clientID {
		t.Errorf(
			"expected client_id %q, got %q",
			clientID,
			query.Get("client_id"),
		)
	}

	if query.Get("scope") != slackBotScopes {
		t.Errorf(
			"expected scope %q, got %q",
			slackBotScopes,
			query.Get("scope"),
		)
	}

	if query.Get("redirect_uri") != redirectURI {
		t.Errorf(
			"expected redirect_uri %q, got %q",
			redirectURI,
			query.Get("redirect_uri"),
		)
	}

	if query.Get("state") != state {
		t.Errorf(
			"expected state %q, got %q",
			state,
			query.Get("state"),
		)
	}
}
