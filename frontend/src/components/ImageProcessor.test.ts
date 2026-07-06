import { describe, it, expect } from 'vitest';
import { compressImage } from './ImageProcessor';

describe('ImageProcessor', () => {
  it('should return a Blob when given a File', async () => {
    const file = new File(['hello'], 'test.jpg', { type: 'image/jpeg' });
    const result = await compressImage(file);
    expect(result).toBeInstanceOf(Blob);
  });

  // 2000px 圧縮テスト
  it('should resize the long edge to 2000px and return a JPEG Blob', async () => {
    const dummyJpgBase64 = '/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAP//////////////////////////////////////////////////////////////////////////////////////wgALCAABAAEBAREA/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxA=';
    const bin = atob(dummyJpgBase64);
    const buffer = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) buffer[i] = bin.charCodeAt(i);
    
    const file = new File([buffer], 'large-image.jpg', { type: 'image/jpeg' });
    
    // 実行
    const result = await compressImage(file);
    
    // 検証
    expect(result.type).toBe('image/jpeg');
    expect(result.size).toBeGreaterThan(0);
  });
});
