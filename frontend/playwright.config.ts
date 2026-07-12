import { defineConfig } from '@playwright/test';

export default defineConfig({
  testMatch: 'tests/**/*.spec.ts',

  use: {
    baseURL: process.env.BASE_URL || 'http://127.0.0.1:8081',
  },

  webServer: {
    command: 'echo "Server already running via Docker"',
    url: 'http://127.0.0.1:8081',
    reuseExistingServer: true,
    timeout: 60 * 1000,
  },
});
