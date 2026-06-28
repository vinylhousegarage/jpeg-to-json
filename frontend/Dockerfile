# 開発用 (Target: dev)
FROM node:24-trixie-slim AS dev
WORKDIR /workspace/frontend
COPY package*.json ./
RUN npm ci
COPY . .
CMD ["npm", "run", "dev"]

# ビルド用 (Target: builder)
FROM node:24-trixie-slim AS builder
WORKDIR /workspace/frontend
COPY package*.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# 本番用 (Target: runner)
FROM node:24-trixie-slim AS runner
WORKDIR /workspace/frontend
ENV NODE_ENV=production
# ビルド成果物のみを builder からコピー
COPY --from=builder /workspace/frontend/.next ./.next
COPY --from=builder /workspace/frontend/public ./public
COPY --from=builder /workspace/frontend/package.json ./package.json
COPY --from=builder /workspace/frontend/node_modules ./node_modules
CMD ["npm", "start"]
