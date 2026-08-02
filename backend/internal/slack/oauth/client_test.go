package oauth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testClientID     = "client-id"
	testClientSecret = "client-secret"
	testCode         = "code123"
	testRedirectURI  = "https://example.com/callback"
)

func TestClient_ExchangeCode_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf(
					"method = %s, want %s",
					r.Method,
					http.MethodPost,
				)
				return
			}

			user, pass, ok := r.BasicAuth()
			if !ok {
				t.Error("missing basic auth")
				return
			}

			if user != testClientID {
				t.Errorf(
					"clientID = %q, want %q",
					user,
					testClientID,
				)
			}

			if pass != testClientSecret {
				t.Errorf(
					"clientSecret = %q, want %q",
					pass,
					testClientSecret,
				)
			}

			if contentType := r.Header.Get("Content-Type"); contentType != "application/x-www-form-urlencoded" {
				t.Errorf(
					"Content-Type = %q, want %q",
					contentType,
					"application/x-www-form-urlencoded",
				)
			}

			if err := r.ParseForm(); err != nil {
				t.Errorf("failed to parse form: %v", err)
				return
			}

			if got := r.Form.Get("code"); got != testCode {
				t.Errorf(
					"code = %q, want %q",
					got,
					testCode,
				)
			}

			if got := r.Form.Get("redirect_uri"); got != testRedirectURI {
				t.Errorf(
					"redirect_uri = %q, want %q",
					got,
					testRedirectURI,
				)
			}

			w.Header().Set("Content-Type", "application/json")

			response := tokenResponse{
				OK:          true,
				AccessToken: "xoxb-test",
				BotUserID:   "B123",
			}
			response.Team.ID = "T123"

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)
	if err != nil {
		t.Fatalf("ExchangeCode() error = %v", err)
	}

	if token.AccessToken != "xoxb-test" {
		t.Errorf(
			"AccessToken = %q, want %q",
			token.AccessToken,
			"xoxb-test",
		)
	}

	if token.BotUserID != "B123" {
		t.Errorf(
			"BotUserID = %q, want %q",
			token.BotUserID,
			"B123",
		)
	}

	if token.TeamID != "T123" {
		t.Errorf(
			"TeamID = %q, want %q",
			token.TeamID,
			"T123",
		)
	}
}

func TestClient_ExchangeCode_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)

	if err == nil {
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf("ExchangeCode() token = %+v, want nil", token)
	}

	if !strings.Contains(err.Error(), "returned status 500") {
		t.Errorf(
			"ExchangeCode() error = %q, want status error",
			err,
		)
	}
}

func TestClient_ExchangeCode_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if _, err := w.Write([]byte("{")); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)

	if err == nil {
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf("ExchangeCode() token = %+v, want nil", token)
	}

	if !strings.Contains(err.Error(), "decode Slack OAuth token response") {
		t.Errorf(
			"ExchangeCode() error = %q, want decode error",
			err,
		)
	}
}

func TestClient_ExchangeCode_SlackError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			response := tokenResponse{
				OK:    false,
				Error: "invalid_code",
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)

	if err == nil {
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf("ExchangeCode() token = %+v, want nil", token)
	}

	if !strings.Contains(err.Error(), "invalid_code") {
		t.Errorf(
			"ExchangeCode() error = %q, want invalid_code",
			err,
		)
	}
}

func TestClient_ExchangeCode_MissingAccessToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			response := tokenResponse{
				OK: true,
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Errorf("failed to encode response: %v", err)
			}
		}),
	)
	defer server.Close()

	client := newTestClient(server)

	token, err := client.ExchangeCode(
		context.Background(),
		testCode,
		testRedirectURI,
	)

	if err == nil {
		t.Fatal("ExchangeCode() error = nil, want an error")
	}

	if token != nil {
		t.Errorf("ExchangeCode() token = %+v, want nil", token)
	}

	if !strings.Contains(err.Error(), "missing access_token") {
		t.Errorf(
			"ExchangeCode() error = %q, want missing access_token error",
			err,
		)
	}
}

func newTestClient(server *httptest.Server) *Client {
	client := NewClient(
		server.Client(),
		testClientID,
		testClientSecret,
	)
	client.tokenURL = server.URL

	return client
}
