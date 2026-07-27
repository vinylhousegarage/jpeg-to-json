import { renderHook, act } from '@testing-library/react';
import { useImageProcessor } from './useImageProcessor';
import { vi, expect } from 'vitest';
import * as imageUtils from '../components/PreviewPhase/utils/compressImage';

vi.spyOn(imageUtils, 'compressImage').mockResolvedValue(new Blob());

describe('useImageProcessor', () => {
  it('compresses image successfully and passes shotNumber', async () => {
    const onCapture = vi.fn();
    const onClearPreview = vi.fn();
    
    const { result } = renderHook(() => useImageProcessor(onCapture, onClearPreview));

    await act(async () => {
      await result.current.processImage(new File([''], 'test.png'));
    });

    expect(result.current.isCompressing).toBe(false);
    expect(onClearPreview).toHaveBeenCalled();
    
    expect(onCapture).toHaveBeenCalledWith(
      expect.any(Blob),
      expect.stringMatching(/^SHOT-\d{3}$/)
    );
  });
});
