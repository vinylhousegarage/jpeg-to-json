import { useAppState } from '../state/useContext';
import { useS3Upload } from '../hooks/useS3Upload';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';
import { Spinner } from '../common/Spinner';

export const Main = () => {
  const { state, dispatch } = useAppState();
  const { upload } = useS3Upload(dispatch);

  const handleSend = async () => {
    if (state.phase.type !== 'preview') {
      return;
    }

    const { file, shotNumber } = state.phase;

    const response = await fetch(
      `${import.meta.env.VITE_API_BASE_URL}/presign`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          filename: `${shotNumber}.jpg`,
          contentType: 'image/jpeg',
          shotNumber,
        }),
      },
    );

    if (!response.ok) {
      dispatch({
        type: 'UPLOAD_COMPLETE',
        status: 'error',
        error: new Error(
          `Failed to get presigned URL: ${response.status}`,
        ),
      });
      return;
    }

    const data: {
      uploadURL: string;
      expiresAt: string;
    } = await response.json();

    await upload(
      file,
      data.uploadURL,
      shotNumber,
    );
  };

  switch (state.phase.type) {
    case 'input':
      return (
        <InputPhase
          isSlackLinked={state.isSlackLinked}
          onConnectSlack={() => {
            dispatch({
              type: 'SET_SLACK_LINKED',
              isSlackLinked: true,
            });
          }}
        />
      );

    case 'preview':
      return (
        <PreviewPhase
          blob={state.phase.file}
          isSending={false}
          onRetake={() => dispatch({ type: 'RETAKE' })}
          onSend={handleSend}
        />
      );

    case 'upload':
      return <Spinner />;

    case 'result':
      return (
        <ResultPhase
          state={state.phase}
          dispatch={dispatch}
        />
      );

    default: {
      const _exhaustiveCheck: never = state.phase;
      return _exhaustiveCheck;
    }
  }
};
