# 共通設定 (Target: base)
FROM node:24-trixie-slim AS base
WORKDIR /workspace/frontend
RUN chown -R node:node /workspace/frontend

# 開発用 (Target: dev)
FROM base AS dev
USER node
COPY --chown=node:node package*.json ./
RUN npm ci
COPY --chown=node:node . .
CMD ["npm", "run", "dev"]

# ビルド用 (Target: builder)
FROM base AS builder
USER node
COPY --chown=node:node package*.json package-lock.json ./
RUN npm ci
COPY --chown=node:node . .
RUN npm run build

# 本番用 (Target: runner)
FROM base AS runner
USER node
ENV NODE_ENV=production
# ビルド成果物のみを builder からコピー
COPY --from=builder --chown=node:node /workspace/frontend/.next ./.next
COPY --from=builder --chown=node:node /workspace/frontend/public ./public
COPY --from=builder --chown=node:node /workspace/frontend/package.json ./package.json
COPY --from=builder --chown=node:node /workspace/frontend/node_modules ./node_modules
CMD ["npm", "start"]
