package storage

import (
	"net/http"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"

	"go.uber.org/zap"
)

// ハンドラー構造体
type PresignHandler struct {
	client *PresignClient
	logger *zap.Logger
}

// コンストラクタ
func NewPresignHandler(client *PresignClient, logger *zap.Logger) *PresignHandler {
	return &PresignHandler{
		client: client,
		logger: logger,
	}
}

func (h *PresignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// リクエストの検証
	req, err := h.validateRequest(r)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	// 署名付きURLの生成
	url, expiresAt, err := h.generatePresignURL(r.Context(), req.Filename)
	if err != nil {
		apierror.WriteError(w, err, h.logger)
		return
	}

	// JSONで返却
	h.writeJSON(w, url, expiresAt)
}
