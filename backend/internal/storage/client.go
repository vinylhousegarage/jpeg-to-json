package storage

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	bucketName  string
	s3Presigner S3Presigner
}

// コンストラクタ
func NewPresignClient(s3Presigner S3Presigner, bucketName string) *PresignClient {
	return &PresignClient{
		bucketName:  bucketName,
		s3Presigner: s3Presigner,
	}
}
