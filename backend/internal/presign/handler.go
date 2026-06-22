package presign

import (
	"context"
	"net/http"
	"time"

	
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"

	"go.uber.org/zap"
)

// 署名付き PutObject リクエスト生成用インターフェース
type S3Presigner interface {
	PresignPutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

// PresignClient 構造体
type PresignClient struct {
	s3Presigner S3Presigner
	bucketName  string
}

// ハンドラー構造体
type PresignHandler struct {
	allowedOrigins []string
	client *PresignClient
	logger *zap.Logger
}

// リクエスト構造体
type PresignRequest struct {
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

// レスポンス構造体
type PresignResponse struct {
  ExpiresAt time.Time `json:"expires_at"`
	UploadURL string    `json:"upload_url"`
}

func (h *PresignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーの設定
	origin := r.Header.Get("Origin")
	h.setCORSHeaders(w, origin)

	// OPTIONSメソッド判定
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
