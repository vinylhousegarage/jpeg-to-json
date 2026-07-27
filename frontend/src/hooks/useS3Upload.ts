import { Dispatch } from 'react';
import { AppAction } from '../types';
import { useAppState } from '../state/AppContext';

export const useS3Upload = (dispatch: Dispatch<AppAction>) => {
  const { state } = useAppState();

  const upload = async (blob: Blob, presignUrl: string) => {
    const slackUrl = state.slackUrl;

    dispatch({ type: 'START_UPLOAD' });
    try {
      const response = await fetch(presignUrl, {
        method: 'PUT',
        body: blob,
        headers: {
          'Content-Type': 'image/jpeg',
          'x-amz-meta-slack-url': encodeURIComponent(slackUrl ?? ''),
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
