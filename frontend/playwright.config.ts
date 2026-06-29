import { defineConfig } from '@playwright/test';

export default defineConfig({
  // テストファイルがどこにあるか
  testDir: './tests', 

  // テストの実行設定
  use: {
    // ブラウザのUIを表示するか（デバッグ用）
    headless: true, 
    // テスト対象のベースURL
    baseURL: 'http://localhost:5173',
  },

  // テスト開始時にサーバーを自動起動する設定
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI, // ローカル環境ならサーバーを再利用する
  },
});
