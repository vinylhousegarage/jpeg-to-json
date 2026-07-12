import { test, expect } from '@playwright/test';

test.describe('S3 Upload Pipeline', () => {
  test('should acquire presigned URL and upload image to S3', async ({ request }) => {
    // ダミーのJPEG画像データを作成 (Base64からBufferへ変換)
    const dummyImageBuffer = Buffer.from(
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=',
      'base64'
    );

    // バックエンドからPresigned URLを取得
    const response = await request.post('/presign', {
      data: { filename: 'test-image.jpg' }
    });

    expect(response.ok()).toBeTruthy();

    // レスポンスから uploadURL を抽出
    const body = await response.json();
    const { uploadURL } = body;
    expect(uploadURL).toBeDefined();

    // S3へ直接PUTリクエストを送信
    const putResponse = await request.put(uploadURL, {
      data: dummyImageBuffer,
      headers: {
        'Content-Type': 'image/jpeg'
      }
    });

    // アップロード結果を検証
    expect(putResponse.ok()).toBeTruthy();
  });
});
