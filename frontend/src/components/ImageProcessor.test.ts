import { describe, it, expect } from 'vitest';
import { compressImage } from './ImageProcessor';

describe('ImageProcessor', () => {
  it('should return a Blob when given a File', async () => {
    const file = new File(['hello'], 'test.jpg', { type: 'image/jpeg' });
    const result = await compressImage(file);
    expect(result).toBeInstanceOf(Blob);
  });

  
});
