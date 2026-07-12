import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests', 
  use: {
    headless: true, 
    baseURL: process.env.BASE_URL || 'http://127.0.0.1:5173',
  },

  webServer: {
    command: 'npm run dev',
    url: 'http://127.0.0.1:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 60 * 1000,
  },
});
