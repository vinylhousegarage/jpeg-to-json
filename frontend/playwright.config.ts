import { defineConfig } from '@playwright/test';

export default defineConfig({
  testMatch: 'tests/**/*.spec.ts',

  use: {
    baseURL: process.env.BASE_URL || 'http://frontend:5173',
  },

  webServer: {
    command: 'echo "Server already running in container"',
    url: 'http://frontend:5173',
    reuseExistingServer: true,
    timeout: 120 * 1000,
  },
});
