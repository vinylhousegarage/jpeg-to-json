import { useState } from 'react';
import { compressImage } from '../components/PreviewPhase/utils/compressImage';

let counter = 1;

export const useImageProcessor = (
  onCapture: (blob: Blob, shotNumber: string) => void,
  onClearPreview: () => void,
  onError?: (err: Error) => void
) => {
  const [isCompressing, setIsCompressing] = useState(false);

  const processImage = async (file: File) => {
    try {
      setIsCompressing(true);
      onClearPreview();
      const compressedBlob = await compressImage(file);

      const shotNumber = `SHOT-${String(counter).padStart(3, '0')}`;
      counter += 1;

      onCapture(compressedBlob, shotNumber);
    } catch (err) {
      onError?.(err instanceof Error ? err : new Error('圧縮に失敗しました'));
    } finally {
      setIsCompressing(false);
    }
  };

  return { isCompressing, processImage };
};
