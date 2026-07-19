package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/bedrock"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/config"
	"github.com/vinylhousegarage/jpeg-to-json/backend/internal/platform/logger"
)

// ディレクトリパス
const (
	samplesDir = "cmd/bedrock-playground/samples"
	resultsDir = "cmd/bedrock-playground/results"
)

func main() {
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

	// 3. 結果ディレクトリの確保
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		l.Fatal("failed to create results directory", zap.String("dir", resultsDir), zap.Error(err))
	}

	// 4. samples/の画像を列挙
	files, err := os.ReadDir(samplesDir)
	if err != nil {
		l.Fatal("failed to read samples directory", zap.String("dir", samplesDir), zap.Error(err))
	}

	// 5. プロンプトを取得
	promptText, err := bedrock.LoadPrompt()
	if err != nil {
    l.Fatal("failed to load prompt", zap.Error(err))
	}

	for _, file := range files {
		// JPEGのみを処理
		if filepath.Ext(file.Name()) != ".jpg" {
			continue
		}

		imagePath := filepath.Join(samplesDir, file.Name())
		l.Info("Processing image", zap.String("path", imagePath))

		// 画像をBase64に変換
		b64Data, err := EncodeFileToBase64(imagePath)
		if err != nil {
			l.Error("failed to encode file to base64", "error", err, "path", imagePath)
			return fmt.Errorf("encode error: %w", err) 
		}

		requestBody := bedrock.NewRequestBody(promptText, b64Data)

		// 出力ファイル名の作成
		name := file.Name()
		baseName := name[:len(name)-len(filepath.Ext(name))]
		outputPath := filepath.Join(resultsDir, baseName+".json")

		// JSONに変換
		jsonData, err := json.MarshalIndent(requestBody, "", "  ")
		if err != nil {
			l.Error("Error marshaling JSON", zap.String("image", name), zap.Error(err))
			continue
		}

		// ファイルに保存
		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			l.Error("Error saving file", zap.String("path", outputPath), zap.Error(err))
			continue
		}

		// ログに出力
		l.Info("Successfully saved request", zap.String("path", outputPath))
	}
}

// テスト用画像を取得
func EncodeFileToBase64(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return bedrock.EncodeBase64(data), nil
}
