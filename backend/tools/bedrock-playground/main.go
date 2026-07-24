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
	samplesDir := "tools/bedrock-playground/samples"

	l.Info("Reading test images from directory...", zap.String("dir", samplesDir))
	files, err := os.ReadDir(samplesDir)
	if err != nil {
		l.Fatal("failed to read samples directory", zap.Error(err))
	}

	for _, file := range files {
		// ディレクトリや隠しファイルはスキップ
		if file.IsDir() {
			continue
		}

		imagePath := fmt.Sprintf("%s/%s", samplesDir, file.Name())
		l.Info("Processing image...", zap.String("path", imagePath))

		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			l.Error("failed to read image file", zap.String("path", imagePath), zap.Error(err))
			continue
		}

		// 解析
		result, err := bedrockService.ProcessImage(ctx, imgData)
		if err != nil {
			l.Error("ProcessImage failed", zap.String("path", imagePath), zap.Error(err))
			continue
		}

		// 結果をターミナルに出力
		fmt.Printf("=== Success: %s ===\n", file.Name())
		for k, v := range result {
			fmt.Printf("  %s: %+v\n", k, v)
		}
		fmt.Println()
	}
}
