package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock/prompts"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
)

// ディレクトリパス
const (
	samplesDir = "cmd/bedrock-playground/samples"
	resultsDir = "cmd/bedrock-playground/results"
)

func main() {
	ctx := context.Background()

	// 1. 設定の初期化
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 2. ロガーの初期化
	l, err := logger.NewLogger(cfg)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer func() {
		_ = l.Sync()
	}()

	// 3. Bedrockクライアントの初期化
	client, err := bedrock.NewClient(ctx, "anthropic.claude-3-5-sonnet-20241022-v2:0")
	if err != nil {
		l.Fatal("failed to initialize bedrock client", zap.Error(err))
	}

	// 4. 結果ディレクトリの確保
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		l.Fatal("failed to create results directory", zap.String("dir", resultsDir), zap.Error(err))
	}

	// 5. samples/の画像を列挙
	files, err := os.ReadDir(samplesDir)
	if err != nil {
		l.Fatal("failed to read samples directory", zap.String("dir", samplesDir), zap.Error(err))
	}

	// 6. プロンプトを取得
	promptText, err := prompts.LoadPrompt("extractor.txt")
	if err != nil {
		l.Fatal("failed to load prompt", zap.Error(err))
	}

	l.Info("Starting playground execution...")

	for _, file := range files {
		// JPEGの確認
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}

		imagePath := filepath.Join(samplesDir, file.Name())
		l.Info("Processing image", zap.String("path", imagePath))

		// 画像ファイルをバイナリとして読み込み
		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			l.Error("Failed to read image", zap.String("path", imagePath), zap.Error(err))
			continue
		}

		// Bedrockの呼び出し
		result, err := client.Invoke(ctx, imgData, promptText)
		if err != nil {
			l.Error("Bedrock invocation failed", zap.String("file", file.Name()), zap.Error(err))
			continue
		}

		// 出力ファイル名の作成
		name := file.Name()
		baseName := name[:len(name)-len(filepath.Ext(name))]
		outputPath := filepath.Join(resultsDir, baseName+".json")

		// 解析結果（map[string]string）を整形して保存
		jsonData, err := json.MarshalIndent(result.Data, "", "  ")
		if err != nil {
			l.Error("Error marshaling response JSON", zap.String("image", name), zap.Error(err))
			continue
		}

		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			l.Error("Error saving file", zap.String("path", outputPath), zap.Error(err))
			continue
		}

		l.Info("Successfully analyzed and saved result", zap.String("path", outputPath))
	}

	l.Info("Playground analysis completed.")
}
