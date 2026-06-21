package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

func setCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

type PresignRequest struct {
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

func validateRequest(r *http.Request) (*PresignRequest, error) {
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

// インターフェース
type S3Presigner interface {
	PresignPutObject(
    ctx context.Context,
    params *s3.PutObjectInput,
    optFns ...func(*s3.PresignOptions),
  ) (*v4.PresignedHTTPRequest, error)
}

// S3 プレサイン用インターフェース
type PresignClient struct {
	s3Presigner S3Presigner
	bucketName  string
}

// PresignClient 生成
func newPresignClient(
  presigner S3Presigner,
  bucketName string,
) *PresignClient {
	return &PresignClient{
		s3Presigner: presigner,
		bucketName:  bucketName,
	}
}

// 署名付きURLを作成
func (c *PresignClient) generatePresignURL(
    ctx context.Context,
    filename string,
) (string, time.Time, error) {
    // 有効期限の起点を設定
    now := time.Now()
    duration := 15 * time.Minute

    // 署名付きURLを生成
    request, err := c.s3Presigner.PresignPutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(c.bucketName),
        Key:    aws.String(filename),
    }, s3.WithPresignExpires(duration))

    if err != nil {
        return "", time.Time{}, fmt.Errorf("failed to sign request: %w", err)
    }

    // 有効期限を計算
    expiresAt := now.Add(duration)

    return request.URL, expiresAt, nil
}

// レスポンス構造体
type PresignResponse struct {
  ExpiresAt time.Time `json:"expires_at"`
	UploadURL string    `json:"upload_url"`
}

func writeJSON(w http.ResponseWriter, url string, expiresAt time.Time) {
	// Content-Type を指定
	w.Header().Set("Content-Type", "application/json")
	// ステータスコードを明示
	w.WriteHeader(http.StatusOK)

	// レスポンスを生成
	resp := PresignResponse{
			ExpiresAt: expiresAt,
			UploadURL: url,
	}

	// エンコード
	if err := json.NewEncoder(w).Encode(resp); err != nil {
			fmt.Printf("failed to encode json: %v\n", err)
	}
}
