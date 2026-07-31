import { test, expect } from '@playwright/test';
import { Locator } from '@playwright/test';

type BoundingBox = NonNullable<Awaited<ReturnType<Locator['boundingBox']>>>;

test.describe('Slack OAuth Phase', () => {
  // iPhone SE2相当のサイズ設定
  test.use({ viewport: { width: 375, height: 667 } });

  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the correct labels and buttons when not linked', async ({ page }) => {
    // 見出しと「Add to Slack」ボタンの存在・テキストを検証
    await expect(page.locator('h2')).toHaveText('Slack通知設定');
    await expect(page.locator('button[type="button"]')).toHaveText('Add to Slack');
  });

  test('should be centered and stacked vertically', async ({ page }) => {
    const heading = page.locator('h2');
    const description = page.locator('p').filter({ hasText: 'Slackアカウントを連携します' });
    const button = page.locator('button[type="button"]');

    const [headingBox, descBox, buttonBox] = await Promise.all([
      heading.boundingBox(),
      description.boundingBox(),
      button.boundingBox()
    ]);

    if (!headingBox || !descBox || !buttonBox) {
      throw new Error('Some elements are not visible');
    }

    const getCenter = (box: BoundingBox) => box.x + box.width / 2;
    
    const headingCenter = getCenter(headingBox);
    const descCenter = getCenter(descBox);
    const buttonCenter = getCenter(buttonBox);

    // 中央寄せの検証（許容誤差1px以内）
    expect(Math.abs(headingCenter - descCenter)).toBeLessThan(1);
    expect(Math.abs(descCenter - buttonCenter)).toBeLessThan(1);

    // 縦並びの検証（上の要素の底辺が下の要素の上辺より上にあるか）
    expect(headingBox.y + headingBox.height).toBeLessThanOrEqual(descBox.y);
    expect(descBox.y + descBox.height).toBeLessThanOrEqual(buttonBox.y);
  });
});
