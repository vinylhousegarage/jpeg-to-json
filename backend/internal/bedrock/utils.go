package bedrock

import "encoding/base64"

// []byteからBase64へ変換
func EncodeBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}
