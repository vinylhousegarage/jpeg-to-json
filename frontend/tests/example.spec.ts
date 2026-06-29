import { test, expect } from '@playwright/test';
import path from 'node:path';

test('has title', async ({ page }) => {
  await page.goto('file://' + path.resolve(__dirname, '../index.html'));
  await expect(page).toHaveTitle(/Playwright/);
});
