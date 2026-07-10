import { test, expect, Page } from '@playwright/test';

test('URL registration form is displayed vertically', async ({ page }: { page: Page }) => {
  await page.goto('/');

  const label = page.locator('label[for="url"]');
  const input = page.locator('input#url');
  const button = page.locator('button[type="submit"]');

  // それぞれの左端の座標を取得
  const labelBox = await label.boundingBox();
  const inputBox = await input.boundingBox();
  const buttonBox = await button.boundingBox();

  // 全ての要素の左端(x)が一致していることを検証
  expect(labelBox?.x).toBe(inputBox?.x);
  expect(inputBox?.x).toBe(buttonBox?.x);
});
