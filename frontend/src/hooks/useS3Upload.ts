import { Dispatch } from 'react';
import { AppAction } from '../types';

export const useS3Upload = (dispatch: Dispatch<AppAction>) => {
  const upload = async (blob: Blob, presignUrl: string, shotNumber: string) => {
    dispatch({ type: 'START_UPLOAD' });
    try {
      const response = await fetch(presignUrl, {
        method: 'PUT',
        body: blob,
        headers: {
          'Content-Type': 'image/jpeg',
          'x-amz-meta-shot-number': encodeURIComponent(shotNumber),
        },
      });

      if (!response.ok) throw new Error('Failed to upload file');

      dispatch({ type: 'UPLOAD_COMPLETE', status: 'success' });
    } catch (e) {
      dispatch({ type: 'UPLOAD_COMPLETE', status: 'error', error: e as Error });
    }
  };

  return { upload };
};
