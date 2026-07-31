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
	logger := zap.NewNop()

	handler := NewHandler(clientID, redirectURI, logger)

	req := httptest.NewRequest(http.MethodGet, "/auth/slack", nil)
	rec := httptest.NewRecorder()

	// 実行
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// ステータスコードの検証
	if res.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, res.StatusCode)
	}

	// Cookieの検証
	cookies := res.Cookies()
	var stateCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "oauth_state" {
			stateCookie = c
			break
		}
	}

	if stateCookie == nil {
		t.Fatal("expected 'oauth_state' cookie to be set, but not found")
	}

	if stateCookie.Value == "" {
		t.Error("expected 'oauth_state' cookie value to not be empty")
	}

	if !stateCookie.HttpOnly || !stateCookie.Secure || stateCookie.Path != "/" {
		t.Errorf("cookie attributes are incorrect: %+v", stateCookie)
	}

	// リダイレクト先URLの検証
	location := res.Header.Get("Location")
	if location == "" {
		t.Fatal("expected Location header for redirection, but got empty")
	}

	parsedURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}

	query := parsedURL.Query()
	if query.Get("client_id") != clientID {
		t.Errorf("expected client_id '%s', got '%s'", clientID, query.Get("client_id"))
	}

	if query.Get("redirect_uri") != redirectURI {
		t.Errorf("expected redirect_uri '%s', got '%s'", redirectURI, query.Get("redirect_uri"))
	}

	if query.Get("state") != stateCookie.Value {
		t.Errorf("expected state in URL to match cookie value, got URL state: '%s', cookie state: '%s'", query.Get("state"), stateCookie.Value)
	}

	if !strings.HasPrefix(location, "https://slack.com/oauth/v2/authorize") {
		t.Errorf("unexpected redirect base URL: %s", location)
	}
}
