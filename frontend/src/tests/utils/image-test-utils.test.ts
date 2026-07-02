import { describe, it, expect } from 'vitest';
import { getImageDimensions } from './image-test-utils';

describe('getImageDimensions', () => {
  it('should correctly return dimensions of a generated image', async () => {
    // Arrange: テスト用のCanvasを作成
    const width = 800;
    const height = 600;
    const canvas = document.createElement('canvas');
    canvas.width = width;
    canvas.height = height;

    // Act: CanvasからBlobを作成し、関数を実行
    const blob = await new Promise<Blob | null>((resolve) => 
      canvas.toBlob(resolve, 'image/png')
    );
    
    if (!blob) throw new Error('Failed to create blob');
    
    const dimensions = await getImageDimensions(blob);

    // Assert: サイズが一致することを確認
    expect(dimensions.width).toBe(width);
    expect(dimensions.height).toBe(height);
  });

  it('should throw an error for invalid blob data', async () => {
    // 空のBlob（画像として読み込めないデータ）を渡してエラーになるか確認
    const invalidBlob = new Blob(['not an image'], { type: 'image/png' });
    
    await expect(getImageDimensions(invalidBlob)).rejects.toThrow(
      "Failed to load image for dimension check"
    );
  });
});
