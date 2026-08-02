package callback

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

func TestHandler_ServeHTTP_RejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	exchanger := &stubCodeExchanger{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		true,
		exchanger,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/oauth/slack/callback",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusMethodNotAllowed,
		apierror.ErrorCodeInvalidMethod,
	)

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Error("unexpected response cookie")
	}
}

func TestHandler_ServeHTTP_MissingStateCookie(t *testing.T) {
	t.Parallel()

	exchanger := &stubCodeExchanger{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		true,
		exchanger,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/slack/callback?state="+testState+"&code="+testCode,
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusBadRequest,
		apierror.ErrorCodeMissingState,
	)

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Error("unexpected response cookie")
	}
}

func TestHandler_ServeHTTP_InvalidState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		queryState  string
		cookieState string
	}{
		{
			name:        "missing query state",
			queryState:  "",
			cookieState: testState,
		},
		{
			name:        "state mismatch",
			queryState:  "different-state",
			cookieState: testState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exchanger := &stubCodeExchanger{}
			store := &stubTokenStore{}

			handler := NewHandler(
				testRedirectURI,
				true,
				exchanger,
				store,
				zap.NewNop(),
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/oauth/slack/callback?state="+tt.queryState+
					"&code="+testCode,
				nil,
			)
			req.AddCookie(&http.Cookie{
				Name:  oauthStateCookieName,
				Value: tt.cookieState,
			})

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assertErrorResponse(
				t,
				rec,
				http.StatusBadRequest,
				apierror.ErrorCodeInvalidState,
			)

			if exchanger.called {
				t.Error("ExchangeCode() was called")
			}

			if store.called {
				t.Error("Save() was called")
			}

			if len(rec.Result().Cookies()) != 0 {
				t.Error(
					"state cookie was deleted before successful validation",
				)
			}
		})
	}
}

func TestHandler_ServeHTTP_MissingCode(t *testing.T) {
	t.Parallel()

	exchanger := &stubCodeExchanger{}
	store := &stubTokenStore{}

	handler := NewHandler(
		testRedirectURI,
		false,
		exchanger,
		store,
		zap.NewNop(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/oauth/slack/callback?state="+testState,
		nil,
	)
	req.AddCookie(&http.Cookie{
		Name:  oauthStateCookieName,
		Value: testState,
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assertErrorResponse(
		t,
		rec,
		http.StatusBadRequest,
		apierror.ErrorCodeMissingCode,
	)

	if exchanger.called {
		t.Error("ExchangeCode() was called")
	}

	if store.called {
		t.Error("Save() was called")
	}

	assertDeleteStateCookie(t, rec, false)
}
