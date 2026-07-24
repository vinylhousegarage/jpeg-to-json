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

// 構造体を定義
type PresignClient struct {
	bucketName  string
	s3Presigner S3Presigner
}

// 構造体を初期化
func NewPresignClient(s3Presigner S3Presigner, bucketName string) *PresignClient {
	return &PresignClient{
		bucketName:  bucketName,
		s3Presigner: s3Presigner,
	}
}

// BucketName Getter
func (c *PresignClient) BucketName() string {
	return c.bucketName
}

// S3Presigner Getter
func (c *PresignClient) S3Presigner() S3Presigner {
	return c.s3Presigner
}
