package get

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// インターフェースを定義
type S3Presigner interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

// 構構造を定義
type Service struct {
	s3Presigner S3Presigner
}

// 構造体を初期化
func NewService(s3Presigner S3Presigner) *Service {
	return &Service{s3Presigner: s3Presigner}
}

// 署名付きURL生成ロジック
func (s *Service) GeneratePresignURL(ctx context.Context, filename string) (string, time.Time, error) {
	now := time.Now()
	duration := 15 * time.Minute

	request, err := s.s3Presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Key: aws.String(filename),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign get request: %w", err)
	}

	return request.URL, now.Add(duration), nil
}
