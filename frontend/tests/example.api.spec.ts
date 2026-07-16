import { test, expect } from '@playwright/test';

test('returns the expected response', async ({ request }) => {
  // レスポンスを取得
  const response = await request.get('/api/example');
  // ステータスコードの確認
  expect(response.status()).toBe(200);

  // ヘッダーを取得
  const headers = response.headers();
  // ヘッダーを確認
  expect(headers['content-type']).toMatch(/application\/json/);

  // ボディを取得
  const body = await response.json();
  // ボディを確認
  expect(body).toEqual({
    id: expect.any(Number),
    message: 'success'
  });
  
  // ログに出力
  console.log('Response headers:', headers);
  console.log('Response body:', body);
});
