import { test, expect } from '@playwright/test';
import { Locator } from '@playwright/test';

type BoundingBox = NonNullable<Awaited<ReturnType<Locator['boundingBox']>>>;

test.describe('URL Registration Form', () => {
  // iPhone SE2相当のサイズ設定
  test.use({ viewport: { width: 375, height: 667 } });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the correct labels and buttons', async ({ page }) => {
    await expect(page.locator('label[for="url"]')).toHaveText('URLを入力');
    await expect(page.locator('button[type="submit"]')).toHaveText('登録して撮影');
  });

  test('should be centered and stacked vertically', async ({ page }) => {
    const label = page.locator('label[for="url"]');
    const input = page.locator('input#url');
    const button = page.locator('button[type="submit"]');

    const [labelBox, inputBox, buttonBox] = await Promise.all([
      label.boundingBox(),
      input.boundingBox(),
      button.boundingBox()
    ]);

    if (!labelBox || !inputBox || !buttonBox) {
      throw new Error('Some elements are not visible');
    }

    // 中央座標の計算
    const getCenter = (box: BoundingBox) => box.x + box.width / 2;
    
    const labelCenter = getCenter(labelBox);
    const inputCenter = getCenter(inputBox);
    const buttonCenter = getCenter(buttonBox);

    // 1. 中央寄せの検証（許容誤差1px以内）
    expect(Math.abs(labelCenter - inputCenter)).toBeLessThan(1);
    expect(Math.abs(inputCenter - buttonCenter)).toBeLessThan(1);

    // 2. 縦並びの検証（上の要素の底辺が、下の要素の上辺より上にあるか）
    // ※ 厳密な重なりチェックのため、各要素の y と y + height を比較
    expect(labelBox.y + labelBox.height).toBeLessThanOrEqual(inputBox.y);
    expect(inputBox.y + inputBox.height).toBeLessThanOrEqual(buttonBox.y);
  });
});
