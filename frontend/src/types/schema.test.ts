import { describe, it, expect } from 'vitest';
import { uploadUrlSchema } from './schema';

describe('uploadUrlSchema', () => {
  it('should succeed with a valid URL', () => {
    const validUrl = 'https://example.com/upload';
    const result = uploadUrlSchema.safeParse(validUrl);
    expect(result.success).toBe(true);
  });

  it('should fail when the URL is empty', () => {
    const emptyUrl = '';
    const result = uploadUrlSchema.safeParse(emptyUrl);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("URLを入力してください");
    }
  });

  it('should fail when the URL format is invalid', () => {
    const invalidUrl = 'not-a-url';
    const result = uploadUrlSchema.safeParse(invalidUrl);
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("有効なURL形式で入力してください");
    }
  });
});
