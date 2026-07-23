package storage

import "time"

// リクエスト構造体
type PresignRequest struct {
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

// レスポンス構造体
type PresignResponse struct {
  ExpiresAt time.Time `json:"expiresAt"`
	UploadURL string    `json:"uploadURL"`
}
