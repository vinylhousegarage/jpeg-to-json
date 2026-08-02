package config

import (
	"strings"
	"testing"
)

func TestLoadSlackConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Setenv("SLACK_CLIENT_ID", "test-client-id")
		t.Setenv("SLACK_CLIENT_SECRET", "test-client-secret")
		t.Setenv(
			"SLACK_REDIRECT_URI",
			"https://example.com/oauth/slack/callback",
		)

		cfg, err := loadSlackConfig()
		if err != nil {
			t.Fatalf("loadSlackConfig() error = %v", err)
		}

		if cfg.ClientID != "test-client-id" {
			t.Errorf(
				"expected ClientID %q, got %q",
				"test-client-id",
				cfg.ClientID,
			)
		}

		if cfg.ClientSecret != "test-client-secret" {
			t.Errorf(
				"expected ClientSecret %q, got %q",
				"test-client-secret",
				cfg.ClientSecret,
			)
		}

		expectedRedirectURI :=
			"https://example.com/oauth/slack/callback"

		if cfg.RedirectURI != expectedRedirectURI {
			t.Errorf(
				"expected RedirectURI %q, got %q",
				expectedRedirectURI,
				cfg.RedirectURI,
			)
		}
	})

	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		redirectURI  string
		expectedErr  string
	}{
		{
			name:         "missing client ID",
			clientID:     "",
			clientSecret: "test-client-secret",
			redirectURI:  "https://example.com/oauth/slack/callback",
			expectedErr:  "SLACK_CLIENT_ID is required",
		},
		{
			name:         "missing client secret",
			clientID:     "test-client-id",
			clientSecret: "",
			redirectURI:  "https://example.com/oauth/slack/callback",
			expectedErr:  "SLACK_CLIENT_SECRET is required",
		},
		{
			name:         "missing redirect URI",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			redirectURI:  "",
			expectedErr:  "SLACK_REDIRECT_URI is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SLACK_CLIENT_ID", tt.clientID)
			t.Setenv("SLACK_CLIENT_SECRET", tt.clientSecret)
			t.Setenv("SLACK_REDIRECT_URI", tt.redirectURI)

			_, err := loadSlackConfig()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf(
					"expected error containing %q, got %q",
					tt.expectedErr,
					err.Error(),
				)
			}
		})
	}
}
