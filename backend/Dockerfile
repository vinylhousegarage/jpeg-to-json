# 開発用 (Target: dev)
FROM golang:1.26.3-alpine3.23 AS dev

# UID/GIDの受け取り
ARG UID=1000
ARG GID=1000

# 依存関係をインストール
RUN apk add --no-cache git build-base zsh curl

# ユーザー作成
RUN addgroup -g ${GID} -S appgroup && \
    adduser -u ${UID} -S appuser -G appgroup -s /bin/zsh

WORKDIR /workspace/backend

# モジュールキャッシュの権限設定
RUN chown appuser:appgroup /workspace/backend /go/pkg/mod /go/bin

# ユーザー切り替え
USER appuser

# パス設定
ENV PATH=$PATH:/go/bin

# Air インストール
RUN go install github.com/air-verse/air@latest

# モジュールダウンロード
COPY --chown=appuser:appgroup go.mod go.sum* ./
RUN go mod download

# ソースコードをコピー
COPY --chown=appuser:appgroup . .

CMD ["air"]


# 本番ビルド用 (Target: builder)
# FROM golang:1.26.3-alpine3.23 AS builder

# ディレクトリ
# WORKDIR /workspace/backend

# 依存関係をキャッシュしてダウンロード
# COPY go.mod go.sum ./
# RUN go mod download

# ソースコードをコピー
# COPY . .

# Lambda (provided.al2023) 用に最適化してビルド
# CGO_ENABLED=0 で静的リンク、ファイル名は 'bootstrap'
# -trimpath と -ldflags="-s -w" でバイナリサイズを削減

# PresignHandler ビルド（bin/presign/）
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#     go build -tags lambda.norpc -trimpath -ldflags="-s -w" \
#     -o /workspace/backend/bin/presign/bootstrap ./cmd/presign/main.go

# MainHandler ビルド（bin/textract/）
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
#     go build -tags lambda.norpc -trimpath -ldflags="-s -w" \
#     -o /workspace/backend/bin/textract/bootstrap ./cmd/textract/main.go
