package apierror

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// ErrorResponse 構造体
type ErrorResponse struct {
	Error string `json:"error"`
}

// エラーを判定してレスポンスを送信
func WriteError(w http.ResponseWriter, err error, logger *zap.Logger) {
	var apiErr *APIError
	
	// デフォルトを status 500 に設定
	status := http.StatusInternalServerError
	code := ErrorCodeInternal
	
	if errors.As(err, &apiErr) {
		status = apiErr.HTTPStatus
		code = apiErr.Code
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	// 構造体を使ってレスポンスを生成
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Error: string(code),
	}); err != nil {
		logger.Error("failed to write json response", zap.Error(err))
	}
}
