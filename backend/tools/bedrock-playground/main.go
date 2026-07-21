package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock/prompts"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
)

func main() {
	// 設定の初期化
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// logger の初期化
	l, err := logger.NewLogger(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer func() {
		_ = l.Sync()
	}()

	ctx := context.Background()

	// プロンプトの読み込み
	promptText, err := prompts.LoadPrompt(cfg.PromptFileName)
	if err != nil {
		l.Fatal("failed to load prompt", zap.Error(err))
	}

	// Bedrock クライアントおよびサービスの初期化
	bedrockClient, err := bedrock.NewClient(ctx, cfg.BedrockModelID, promptText, l)
	if err != nil {
		l.Fatal("failed to create bedrock client", zap.Error(err))
	}
	bedrockService := bedrock.NewService(bedrockClient)

	// テスト用画像の読み込み
	imagePath := "sample.jpg"
	l.Info("Reading test image...", zap.String("path", imagePath))
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		l.Fatal("failed to read test image (make sure sample.jpg exists)", zap.Error(err))
	}

	// 解析
	l.Info("Starting Bedrock inference...")
	result, err := bedrockService.ProcessImage(ctx, imgData)
	if err != nil {
		l.Fatal("ProcessImage failed", zap.Error(err))
	}

	// 結果をターミナルに出力
	fmt.Println("Bedrock Analysis Success!")
	for k, v := range result {
		fmt.Printf("  %s: %+v\n", k, v)
	}
}
