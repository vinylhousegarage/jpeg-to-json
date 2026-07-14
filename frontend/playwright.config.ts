import { defineConfig } from '@playwright/test';

export default defineConfig({
  testMatch: 'tests/**/*.spec.ts',

  use: {
    baseURL: process.env.BASE_URL || 'http://127.0.0.1:5173',
  },

  webServer: {
    command: 'sleep infinity',
    url: 'http://127.0.0.1:5173',
    reuseExistingServer: true,
    timeout: 60 * 1000,
  },
});
