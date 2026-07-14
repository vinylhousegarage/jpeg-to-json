import { test, expect } from '@playwright/test';

test.describe('S3 Upload Pipeline', () => {
  test('should acquire presigned URL and upload image to S3', async ({ request }) => {
    // ダミーのJPEG画像データを作成 (Base64 から Buffer へ変換)
    const dummyImageBuffer = Buffer.from(
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=',
      'base64'
    );

    // バックエンドから Presigned URL を取得
    const API_BASE_URL = process.env.API_BASE_URL || 'http://backend:8080' as string;
    const response = await request.post(`${ API_BASE_URL }/presign`, {
      headers: { 'Content-Type': 'application/json' },
      data: { filename: 'test-image.jpg' }
    });

    // エラー発生時に詳細を表示
    if (!response.ok()) {
      const errorBody = await response.text();
      console.error(`[DEBUG] Presigned URL acquisition failed (Status: ${response.status()})`);
      console.error(`[DEBUG] Response Body: ${errorBody}`);
      
      // エラー内容
      expect(response.ok(), `Failed to acquire presigned URL: ${response.status()} - ${errorBody}`).toBeTruthy();
    }

    // レスポンスから uploadURL を抽出
    const body = await response.json();
    const { uploadURL } = body;
    expect(uploadURL).toBeDefined();

    // S3 へ直接 PUT リクエストを送信
    const putResponse = await request.put(uploadURL, {
      data: dummyImageBuffer,
      headers: {
        'Content-Type': 'image/jpeg'
      }
    });

    // アップロード結果の検証とエラー出力
    if (!putResponse.ok()) {
      const errorBody = await putResponse.text();
      console.error(`[DEBUG] S3 Upload failed (Status: ${putResponse.status()})`);
      console.error(`[DEBUG] Response Body: ${errorBody}`);
      
      expect(putResponse.ok(), `Failed to upload to S3: ${putResponse.status()} - ${errorBody}`).toBeTruthy();
    }
  });
});
