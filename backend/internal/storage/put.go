package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

// リクエスト検証
func (h *PresignHandler) validatePutPresignRequest(r *http.Request) (*PutPresignRequest, error) {
	if r.Method != http.MethodPost {
		return nil, apierror.New(apierror.ErrorCodeInvalidMethod, http.StatusMethodNotAllowed, nil)
	}

	var req PutPresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierror.New(apierror.ErrorCodeInvalidJSON, http.StatusBadRequest, err)
	}

	if req.Filename == "" {
		return nil, apierror.New(apierror.ErrorCodeMissingFilename, http.StatusBadRequest, nil)
	}

	return &req, nil
}

// PutPresignURL 生成
func (h *PresignHandler) generatePutPresignURL(ctx context.Context, filename string) (string, time.Time, error) {
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
func (h *PresignHandler) writeUploadResponse(w http.ResponseWriter, url string, expiresAt time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := PutPresignResponse{
		ExpiresAt: expiresAt,
		UploadURL: url,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode json", zap.Error(err))
	}
}
