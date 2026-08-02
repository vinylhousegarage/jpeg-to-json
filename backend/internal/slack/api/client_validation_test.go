package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestClient_PostMessage_ValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accessToken string
		channelID   string
		text        string
		wantError   string
	}{
		{
			name:        "missing access token",
			accessToken: "",
			channelID:   testChannelID,
			text:        testMessageText,
			wantError:   "post slack message: access token is empty",
		},
		{
			name:        "missing channel ID",
			accessToken: testAccessToken,
			channelID:   "",
			text:        testMessageText,
			wantError:   "post slack message: channel ID is empty",
		},
		{
			name:        "missing text",
			accessToken: testAccessToken,
			channelID:   testChannelID,
			text:        "",
			wantError:   "post slack message: text is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transportCalled := false

			httpClient := &http.Client{
				Transport: roundTripFunc(
					func(*http.Request) (*http.Response, error) {
						transportCalled = true

						return nil, errors.New("unexpected request")
					},
				),
			}

			client := NewClient(httpClient)

			err := client.PostMessage(
				context.Background(),
				tt.accessToken,
				tt.channelID,
				tt.text,
			)
			if err == nil {
				t.Fatal("PostMessage() error = nil, want an error")
			}

			if err.Error() != tt.wantError {
				t.Errorf(
					"PostMessage() error = %q, want %q",
					err.Error(),
					tt.wantError,
				)
			}

			if transportCalled {
				t.Error("HTTP request was sent for invalid input")
			}
		})
	}
}
