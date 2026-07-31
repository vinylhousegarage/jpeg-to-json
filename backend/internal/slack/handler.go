package slack

import (
	"net/http"

	"go.uber.org/zap"
)

// 構造体を定義
type Handler struct {
	clientID    string
	redirectURI string
	logger      *zap.Logger
}

// 構造体を初期化
func NewHandler(clientID, redirectURI string, logger *zap.Logger) *Handler {
	return &Handler{
		clientID:    clientID,
		redirectURI: redirectURI,
		logger:      logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := GenerateState()
	cookie := BuildStateCookie(state)
	http.SetCookie(w, cookie)

	authURL := BuildAuthURL(h.clientID, h.redirectURI, state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}
