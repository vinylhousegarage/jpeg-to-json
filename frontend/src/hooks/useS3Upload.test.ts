import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useS3Upload } from './useS3Upload';

describe('useS3Upload', () => {
  const mockDispatch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
    }));
  });

  it('dispatches correctly on successful completion', async () => {
    const { upload } = useS3Upload(mockDispatch);
    const mockBlob = new Blob(['test-image'], { type: 'image/jpeg' });
    
    await upload(mockBlob, 'https://test.url');

    expect(mockDispatch).toHaveBeenCalledWith({ type: 'START_UPLOAD' });
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'UPLOAD_COMPLETE',
      status: 'success',
    });
  });
});
