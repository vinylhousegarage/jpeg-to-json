package storage

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockPresigner struct{}

func (m *MockPresigner) PresignPutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {

	return &v4.PresignedHTTPRequest{
		URL: "https://example.com/test.jpg",
	}, nil
}
