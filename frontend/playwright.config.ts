import { defineConfig } from '@playwright/test';

export default defineConfig({
  use: {
    baseURL: process.env.BASE_URL || 'http://127.0.0.1:8081',
  },
  webServer: {
    // 起動コマンドは compose.yaml
    command: 'echo "Server already running via Docker"',
    url: 'http://127.0.0.1:8081',
    // compose.yaml 起動のサーバーに接続
    reuseExistingServer: true,
    // サーバーに接続するまで最大60秒待機
    timeout: 60 * 1000,
  },
});
