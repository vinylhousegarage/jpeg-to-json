package slack

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	clientID     string
	redirectURI  string
	cookieSecure bool
	logger       *zap.Logger
}

func NewHandler(
	clientID string,
	redirectURI string,
	cookieSecure bool,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		clientID:     clientID,
		redirectURI:  redirectURI,
		cookieSecure: cookieSecure,
		logger:       logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(
			w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}

	state, err := GenerateState()
	if err != nil {
		h.logger.Error(
			"failed to generate Slack OAuth state",
			zap.Error(err),
		)

		http.Error(
			w,
			"failed to start Slack OAuth",
			http.StatusInternalServerError,
		)
		return
	}

	http.SetCookie(
		w,
		BuildStateCookie(state, h.cookieSecure),
	)

	authURL := BuildAuthURL(
		h.clientID,
		h.redirectURI,
		state,
	)

	http.Redirect(
		w,
		r,
		authURL,
		http.StatusFound,
	)
}
