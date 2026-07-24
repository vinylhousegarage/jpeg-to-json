package storage

import "time"

// Put用リクエスト構造体
type PutPresignRequest struct {
	Filename string `json:"filename"`
	ContentType string `json:"contentType"`
}

// Put用レスポンス構造体
type PutPresignResponse struct {
	ExpiresAt time.Time `json:"expiresAt"`
	UploadURL string    `json:"uploadURL"`
}
