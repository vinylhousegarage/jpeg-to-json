import { describe, it, expect } from 'vitest';
import { compressImage } from './compressImage';
import { getImageDimensions } from './getImageDimensions';

describe('compressImage', () => {
  // 変換をテスト
  it('should return a Blob when given a File', async () => {
    const file = new File(['hello'], 'test.jpg', { type: 'image/jpeg' });
    const result = await compressImage(file);
    
    expect(result).toBeInstanceOf(Blob);
    expect(result.type).toBe('image/jpeg');
  });

  // 2000px制限の圧縮をテスト
  it('should resize the long edge to 2000px or less and return a JPEG Blob', async () => {
    // 1x1ピクセルの最小限の正当なJPEGダミーデータ
    const dummyJpgBase64 = '/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////wgALCAABAAEBAREA/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxA=';
    const bin = atob(dummyJpgBase64);
    const buffer = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) {
      buffer[i] = bin.charCodeAt(i);
    }

    const file = new File([buffer], 'large-image.jpg', { type: 'image/jpeg' });

    // 圧縮ロジックの実行
    const result = await compressImage(file);

    // 検証1: 出力形式
    expect(result.type).toBe('image/jpeg');
    expect(result.size).toBeGreaterThan(0);

    // 検証2: サイズ
    const dimensions = await getImageDimensions(result);
    expect(dimensions.width).toBeLessThanOrEqual(2000);
    expect(dimensions.height).toBeLessThanOrEqual(2000);
  });
});
