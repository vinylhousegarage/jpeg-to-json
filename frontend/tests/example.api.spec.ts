import { test, expect } from '@playwright/test';

describe('GET /api/example', () => {
  it('returns the expected response', async () => {
    // 実際にリクエストを送信
    const response = await request(app).get('/api/example');

    // 1. ステータスコードの確認
    expect(response.status).toBe(200);

    // 2. レスポンスヘッダーの確認
    expect(response.headers['content-type']).toMatch(/json/);

    // 3. ボディの中身を詳しく確認
    expect(response.body).toEqual({
      id: expect.any(Number),
      message: 'success'
    });
    
    // 中身をログに出したい場合（デバッグ用）
    console.log('Response body:', response.body);
  });
});
