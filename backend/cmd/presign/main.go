package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/presign"
)

func main() {
  // 設定の初期化
  cfg, err := config.LoadConfig()
  if err != nil {
    log.Fatalf("failed to load config: %v", err)
  }

  // AWS 設定の読み込み
  awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.Region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	// S3 クライアント
  baseS3Client := s3.NewFromConfig(awsCfg)
	// 署名専用クライアントに変換
	s3PresignClient := s3.NewPresignClient(baseS3Client)
	// 依存を注入
  presignClient := presign.NewPresignClient(s3PresignClient, cfg.InputBucketName)

  // logger の初期化
  l, err := logger.NewLogger(cfg)
  if err != nil {
    panic(fmt.Sprintf("failed to initialize logger: %v", err))
  }

  // ログに出力
  l.Info("Application successfully initialized")

  // ハンドラーの初期化
  presignHandler := presign.NewPresignHandler(cfg.AllowedOrigins, presignClient, l)

  // サーバーの初期化
  mux := http.NewServeMux()

  // ルーティング
  mux.Handle("/presign", presignHandler)

  // サーバー起動
	if cfg.IsLambda {
		// Lambda 環境
		l.Info("Starting server on AWS Lambda")

		adapter := httpadapter.NewV2(mux)
		lambda.Start(adapter.ProxyWithContext)

	} else {
		// ローカル環境
		l.Info("Starting local server on :8080")

		srv := &http.Server{Addr: ":8080", Handler: mux}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}
}
