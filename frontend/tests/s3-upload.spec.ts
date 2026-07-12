import { test, expect } from '@playwright/test';

test.describe('S3 Upload Pipeline', () => {
  test('should acquire presigned URL and upload image to S3', async ({ request }) => {
    // ダミーのJPEG画像データを作成 (Base64からBufferへ変換)
    const dummyImageBuffer = Buffer.from(
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=',
      'base64'
    );

    // バックエンドからPresigned URLを取得
    const response = await request.post('http://backend:8080/presign', {
      data: { filename: 'test-image.jpg' }
    });

    // --- デバッグ用出力 ---
    if (!response.ok()) {
      console.log('--- [DEBUG] Presigned URL acquisition failed ---');
      console.log('Status:', response.status());
      console.log('Body:', await response.text());
    }
    // -----------------------

    expect(response.ok()).toBeTruthy();

    // レスポンスから uploadURL を抽出
    const body = await response.json();
    const { uploadURL } = body;
    expect(uploadURL).toBeDefined();

    // 2. S3へ直接PUTリクエストを送信
    const putResponse = await request.put(uploadURL, {
      data: dummyImageBuffer,
      headers: {
        'Content-Type': 'image/jpeg'
      }
    });

    // --- デバッグ用出力 ---
    if (!putResponse.ok()) {
      console.log('--- [DEBUG] S3 Upload failed ---');
      console.log('Status:', putResponse.status());
      console.log('Body:', await putResponse.text());
    }
    // -----------------------

    // アップロード結果を検証
    expect(putResponse.ok()).toBeTruthy();
  });
});
