package callback

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/slack/oauth"
)

func TestHandler_ServeHTTP_ExchangeCodeError(t *testing.T) {
	t.Parallel()

	exchangeErr := errors.New("slack token exchange failed")
	exchanger := &stubCodeExchanger{
		err: exchangeErr,
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

	assertErrorResponse(
		t,
		rec,
		http.StatusInternalServerError,
		apierror.ErrorCodeInternal,
	)

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
	}

	if store.called {
		t.Error("Save() was called after ExchangeCode() failed")
	}

	assertDeleteStateCookie(t, rec, true)
}

func TestHandler_ServeHTTP_SaveError(t *testing.T) {
	t.Parallel()

	token := &oauth.Token{
		AccessToken: "xoxb-test",
		BotUserID:   "B123",
		TeamID:      "T123",
	}
	exchanger := &stubCodeExchanger{
		token: token,
	}
	store := &stubTokenStore{
		err: errors.New("dynamodb save failed"),
	}

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

	assertErrorResponse(
		t,
		rec,
		http.StatusInternalServerError,
		apierror.ErrorCodeInternal,
	)

	if !exchanger.called {
		t.Fatal("ExchangeCode() was not called")
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
