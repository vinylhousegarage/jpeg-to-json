package presign

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"

	"go.uber.org/zap"
)

// CORSヘッダー設定
func (h *PresignHandler) setCORSHeaders(w http.ResponseWriter, origin string) {
	isAllowed := false
	for _, allowed := range h.allowedOrigins {
		if allowed == origin {
			isAllowed = true
			break
		}
	}

	if isAllowed {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// リクエスト検証
func (h *PresignHandler) validateRequest(r *http.Request) (*PresignRequest, error) {
	if r.Method != http.MethodPost {
		return nil, apierror.New(apierror.ErrorCodeInvalidMethod, http.StatusMethodNotAllowed, nil)
	}

	var req PresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierror.New(apierror.ErrorCodeInvalidJSON, http.StatusBadRequest, err)
	}

	if req.Filename == "" {
		return nil, apierror.New(apierror.ErrorCodeMissingFilename, http.StatusBadRequest, nil)
	}

	return &req, nil
}

// PresignURL 生成
func (h *PresignHandler) generatePresignURL(ctx context.Context, filename string) (string, time.Time, error) {
	now := time.Now()
	duration := 15 * time.Minute

	request, err := h.client.s3Presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(h.client.bucketName),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign request: %w", err)
	}

	return request.URL, now.Add(duration), nil
}

// JSONレスポンス書き込み
func (h *PresignHandler) writeJSON(w http.ResponseWriter, url string, expiresAt time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := PresignResponse{
		ExpiresAt: expiresAt,
		UploadURL: url,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
