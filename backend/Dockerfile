# 開発用 (Target: dev)
FROM golang:1.26.3-alpine3.23 AS dev
WORKDIR /app

# g基本ツールをインストール
RUN apk add --no-cache git build-base

# Air のインストール
RUN go install github.com/air-verse/air@latest

# モジュールのダウンロード
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# ホットリロードを有効化
CMD ["air"]

# 本番ビルド用 (Target: builder)

# FROM golang:1.26.3-alpine3.23 AS builder
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .

# Lambda (provided.al2023) 用に最適化してビルド
# CGO_ENABLED=0 で静的リンク、ファイル名は 'bootstrap'
# PresignHandler ビルド（bin/presign/）
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#     go build -tags lambda.norpc -o /app/bin/presign/bootstrap ./cmd/presign/main.go

# MainHandler ビルド（bin/textract/）
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#     go build -tags lambda.norpc -o /app/bin/textract/bootstrap ./cmd/textract/main.go
