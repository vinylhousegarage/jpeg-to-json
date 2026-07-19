package bedrock

import (
	"embed"
	"encoding/base64"
	"fmt"
)

// プロンプトファイルを保持するためのグローバル変数
//go:embed prompts/extractor.txt
var PromptFS embed.FS

// 埋め込まれたプロンプトの読み込み
func LoadPrompt() (string, error) {
	data, err := PromptFS.ReadFile("prompts/extractor.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded prompt: %w", err)
	}
	return string(data), nil
}

// []byteからBase64へ変換
func EncodeBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}
