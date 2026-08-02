package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
	t.Parallel()

	token := &oauth.Token{
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		TeamID:      "T123",
	}
	exchanger := &stubCodeExchanger{
		token: token,
	}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		true,
		exchanger,
		store,
		zap.NewNop(),
	)

	req := newValidCallbackRequest()
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			rec.Code,
			http.StatusNoContent,
			rec.Body.String(),
		)
	}

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if exchanger.gotCode != testCode {
		t.Errorf(
			"ExchangeCode() code = %q, want %q",
			exchanger.gotCode,
			testCode,
		)
	}

	if exchanger.gotRedirect != testRedirectURI {
		t.Errorf(
			"ExchangeCode() redirectURI = %q, want %q",
			exchanger.gotRedirect,
			testRedirectURI,
		)
	}

	if !store.called {
		t.Fatal("Save() was not called")
	}

	if store.gotToken != token {
		t.Errorf(
			"Save() token = %+v, want %+v",
			store.gotToken,
			token,
		)
	}

	assertDeleteStateCookie(t, rec, true)
}
