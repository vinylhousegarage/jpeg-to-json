import { describe, it, expect, vi, beforeEach } from 'vitest';
import { compressImage } from './compressImage';
import { getImageDimensions } from './getImageDimensions';

describe('compressImage', () => {
  beforeEach(() => {
    // グローバルオブジェクト（URLなど）のモックが必要な場合はここでリセット
    vi.restoreAllMocks();
  });

  // 変換をテスト
  it('should return a Blob when given a File', async () => {
    const file = new File(['hello'], 'test.jpg', { type: 'image/jpeg' });
    const result = await compressImage(file);
    
    expect(result).toBeInstanceOf(Blob);
    expect(result.type).toBe('image/jpeg');
  });

  // 2000px 制限の圧縮をテスト
  it('should resize the long edge to 2000px or less and return a JPEG Blob', async () => {
    const dummyJpgBase64 = '/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////wgALCAABAAEBAREA/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxA=';
    const bin = atob(dummyJpgBase64);
    const buffer = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) {
      buffer[i] = bin.charCodeAt(i);
    }

    const file = new File([buffer], 'large-image.jpg', { type: 'image/jpeg' });
    const result = await compressImage(file);

    expect(result.type).toBe('image/jpeg');
    expect(result.size).toBeGreaterThan(0);

    const dimensions = await getImageDimensions(result);
    expect(dimensions.width).toBeLessThanOrEqual(2000);
    expect(dimensions.height).toBeLessThanOrEqual(2000);
  });

  // 品質 85% 指定をテスト
  it('should call canvas.toBlob with quality parameter set to 0.85', async () => {
    // HTMLCanvasElement のプロトタイプから toBlob を監視
    const toBlobSpy = vi.spyOn(HTMLCanvasElement.prototype, 'toBlob');

    const file = new File(['dummy content'], 'test.jpg', { type: 'image/jpeg' });
    await compressImage(file);

    // toBlob が呼び出された際の引数 (callback, type, quality) を検証
    expect(toBlobSpy).toHaveBeenCalled();
    const mostRecentCall = toBlobSpy.mock.calls[0];
    
    expect(mostRecentCall[1]).toBe('image/jpeg'); // 第2引数がJPEG形式か
    expect(mostRecentCall[2]).toBe(0.85);         // 第3引数が「0.85（品質85%）」か
  });
});
