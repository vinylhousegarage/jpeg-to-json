# ------------------------------
# 共通設定 (Target: base)
# ------------------------------
FROM node:24-trixie-slim AS base

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

# ユーザー切り替え
USER node

# 設定ファイルをコピー
COPY --chown=node:node package*.json ./

# npm をインストール
RUN npm ci

# ソースコードをコピー
COPY --chown=node:node . .

# Vite 開発サーバーのポートを開放
EXPOSE 5173

# --host 0.0.0.0 からのアクセスを許可
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

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
