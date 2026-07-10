import { test, expect, Page } from '@playwright/test';
import { Locator } from '@playwright/test';

type BoundingBox = NonNullable<Awaited<ReturnType<Locator['boundingBox']>>>;

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

    // nullチェック（要素が存在しない可能性を考慮）
    if (!labelBox || !inputBox || !buttonBox) {
      throw new Error('Some elements are not visible');
    }

    // 中央座標の計算（型を BoundingBox と定義）
    const getCenter = (box: BoundingBox) => box.x + box.width / 2;
    
    const labelCenter = getCenter(labelBox);
    const inputCenter = getCenter(inputBox);
    const buttonCenter = getCenter(buttonBox);

    // 全ての要素の中央が一致しているか（中央寄せの検証）
    expect(Math.abs(labelCenter - inputCenter)).toBeLessThan(1);
    expect(Math.abs(inputCenter - buttonCenter)).toBeLessThan(1);

    // 縦並びの検証
    expect(labelBox.y + labelBox.height).toBeLessThan(inputBox.y);
    expect(inputBox.y + inputBox.height).toBeLessThan(buttonBox.y);
  });
});
