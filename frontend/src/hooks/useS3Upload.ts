import { Dispatch } from 'react';
import { AppAction } from '../types';

export const useS3Upload = (dispatch: Dispatch<AppAction>) => {
  const upload = async (blob: Blob, url: string) => {
    dispatch({ type: 'START_UPLOAD' });
    try {
      const response = await fetch(url, {
        method: 'PUT',
        body: blob,
        headers: {
          'Content-Type': 'image/jpeg',
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
