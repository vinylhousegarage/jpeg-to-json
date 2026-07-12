# ------------------------------
# 共通設定 (Target: base)
# ------------------------------
FROM node:24-trixie-slim AS base

# パッケージリストの更新とツールのインストール
RUN apt-get update && apt-get install -y \
    git \
    zsh \
    && rm -rf /var/lib/apt/lists/*

# UID/GIDの受け取り
ARG UID=1000
ARG GID=1000

# ユーザー作成
RUN usermod -u ${UID} node && groupmod -g ${GID} node

# ディリレクトリ
WORKDIR /workspace/frontend

# Zshプラグインのインストール
RUN git clone https://github.com/zsh-users/zsh-autosuggestions /home/node/.oh-my-zsh/custom/plugins/zsh-autosuggestions \
    && echo 'source /home/node/.oh-my-zsh/custom/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh' >> /home/node/.zshrc

# 所有者をユーザーに変更
RUN chown -R node:node /workspace/frontend

# ------------------------------------
# --- 開発用 (Target: dev) ---
# ------------------------------------
FROM base AS dev

# root を明示
USER root

# 依存関係をインストール
RUN apt-get update && apt-get install -y \
    libnss3 \
    libnspr4 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxfixes3 \
    libxrandr2 \
    libgbm1 \
    libasound2t64 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libxcursor1 \
    libgtk-3-0 \
    libgdk-pixbuf-xlib-2.0-0 \
    libcairo-gobject2 \
    libxss1 \
    && rm -rf /var/lib/apt/lists/*

# ChromiumブラウザとOS依存ライブラリをインストール
RUN npx playwright install chromium

# ユーザー切り替え
USER node

# Git 所有権エラー回避
RUN git config --global --add safe.directory '*'

# 設定ファイルをコピー
COPY --chown=node:node package*.json ./

# npm をインストール
RUN npm ci

# ソースコードをコピー
COPY --chown=node:node . .

# Vite 開発サーバーのポートを開放
EXPOSE 5173

# ------------------------------------
# --- ビルド用 (Target: builder) ---
# ------------------------------------
FROM base AS builder
USER node
COPY --chown=node:node package*.json package-lock.json ./
RUN npm install
COPY --chown=node:node . .
# Vite でビルドして dist フォルダを作成
RUN npm run build

# ------------------------------------
# --- 本番用 (Target: runner) ---
# ------------------------------------
FROM nginx:alpine AS runner
# ビルド成果物を Nginx の公開ディリレクトリにコピー
COPY --from=builder /workspace/frontend/dist /usr/share/nginx/html
# Nginx ポート
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
