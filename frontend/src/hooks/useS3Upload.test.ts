import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useS3Upload } from './useS3Upload';
import * as AppContextModule from '../state/AppContext';

vi.mock('../state/AppContext', async () => {
  const actual = await vi.importActual('../state/AppContext');
  return {
    ...actual,
    useAppState: vi.fn(),
  };
});

describe('useS3Upload', () => {
  const mockDispatch = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

  vi.spyOn(AppContextModule, 'useAppState').mockReturnValue({
    state: {
      type: 'preview',
      slackUrl: 'https://hooks.slack.com/services/test',
      shotNumber: 'SHOT-001',
      file: new Blob(['test'], { type: 'image/jpeg' }),
      previewUrl: 'blob:http://localhost/test',
    },
    dispatch: mockDispatch,
  });

    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
    }));
  });

  it('dispatches correctly on successful completion and sends metadata', async () => {
    const { upload } = useS3Upload(mockDispatch);
    const mockBlob = new Blob(['test-image'], { type: 'image/jpeg' });

    await upload(mockBlob, 'https://test-s3-presign.url');

    expect(fetch).toHaveBeenCalledWith(
      'https://test-s3-presign.url',
      expect.objectContaining({
        method: 'PUT',
        headers: expect.objectContaining({
          'x-amz-meta-slack-url': encodeURIComponent('https://hooks.slack.com/services/test'),
        }),
      })
    );

    expect(mockDispatch).toHaveBeenCalledWith({ type: 'START_UPLOAD' });
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'UPLOAD_COMPLETE',
      status: 'success',
    });
  });
});
