package slack

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	clientID := "test-client-id"
	redirectURI := "https://example.com/callback"
	cookieSecure := true
	logger := zap.NewNop()

	handler := NewHandler(
		clientID,
		redirectURI,
		cookieSecure,
		logger,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/slack/start",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusFound {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusFound,
			res.StatusCode,
		)
	}

	cookies := res.Cookies()

	var stateCookie *http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == oauthStateCookieName {
			stateCookie = cookie
			break
		}
	}

	if stateCookie == nil {
		t.Fatalf(
			"expected %q cookie to be set",
			oauthStateCookieName,
		)
	}

	if stateCookie.Value == "" {
		t.Error("expected OAuth state cookie value to not be empty")
	}

	if !stateCookie.HttpOnly {
		t.Error("expected OAuth state cookie HttpOnly to be true")
	}

	if !stateCookie.Secure {
		t.Error("expected OAuth state cookie Secure to be true")
	}

	if stateCookie.Path != oauthStateCookiePath {
		t.Errorf(
			"expected OAuth state cookie Path %q, got %q",
			oauthStateCookiePath,
			stateCookie.Path,
		)
	}

	if stateCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf(
			"expected OAuth state cookie SameSite %v, got %v",
			http.SameSiteLaxMode,
			stateCookie.SameSite,
		)
	}

	expectedMaxAge := int(OAuthStateTTL.Seconds())
	if stateCookie.MaxAge != expectedMaxAge {
		t.Errorf(
			"expected OAuth state cookie MaxAge %d, got %d",
			expectedMaxAge,
			stateCookie.MaxAge,
		)
	}

	if stateCookie.Expires.IsZero() {
		t.Error("expected OAuth state cookie Expires to be set")
	}

	location := res.Header.Get("Location")
	if location == "" {
		t.Fatal("expected Location header, got empty value")
	}

	parsedURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}

	if parsedURL.Scheme != "https" {
		t.Errorf(
			"expected redirect scheme %q, got %q",
			"https",
			parsedURL.Scheme,
		)
	}

	if parsedURL.Host != "slack.com" {
		t.Errorf(
			"expected redirect host %q, got %q",
			"slack.com",
			parsedURL.Host,
		)
	}

	if parsedURL.Path != "/oauth/v2/authorize" {
		t.Errorf(
			"expected redirect path %q, got %q",
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

	if query.Get("state") != stateCookie.Value {
		t.Errorf(
			"expected URL state %q to match cookie state %q",
			query.Get("state"),
			stateCookie.Value,
		)
	}

	if !strings.HasPrefix(
		location,
		slackAuthorizeURL,
	) {
		t.Errorf(
			"expected redirect URL to start with %q, got %q",
			slackAuthorizeURL,
			location,
		)
	}
}

func TestHandler_ServeHTTP_MethodNotAllowed(t *testing.T) {
	t.Parallel()

	handler := NewHandler(
		"test-client-id",
		"https://example.com/callback",
		true,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/oauth/slack/start",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			res.StatusCode,
		)
	}

	if allow := res.Header.Get("Allow"); allow != http.MethodGet {
		t.Errorf(
			"expected Allow header %q, got %q",
			http.MethodGet,
			allow,
		)
	}
}
