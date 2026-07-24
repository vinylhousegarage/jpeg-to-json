package storage

import "time"

// Put用リクエスト構造体
type PutPresignRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

// Put用レスポンス構造体
type PutPresignResponse struct {
	ExpiresAt time.Time `json:"expiresAt"`
	UploadURL string    `json:"uploadURL"`
}

// Get用リクエスト構造体
type GetPresignRequest struct {
	Filename string `json:"filename"`
}

// Get用レスポンス構造体
type GetPresignResponse struct {
	ExpiresAt  time.Time `json:"expiresAt"`
	DownloadURL string   `json:"downloadURL"`
}
