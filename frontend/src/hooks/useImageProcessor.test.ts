import { renderHook, act } from '@testing-library/react';
import { useImageProcessor } from './useImageProcessor';
import { vi } from 'vitest';
import * as imageUtils from '../components/PreviewPhase/utils/compressImage';

// モック定義
vi.spyOn(imageUtils, 'compressImage').mockResolvedValue(new Blob());

describe('useImageProcessor', () => {
  it('compresses image successfully', async () => {
    const onCapture = vi.fn();
    const onClearPreview = vi.fn();
    
    const { result } = renderHook(() => useImageProcessor(onCapture, onClearPreview));

    await act(async () => {
      await result.current.processImage(new File([''], 'test.png'));
    });

    expect(result.current.isCompressing).toBe(false);
    expect(onClearPreview).toHaveBeenCalled();
    expect(onCapture).toHaveBeenCalled();
  });
});
