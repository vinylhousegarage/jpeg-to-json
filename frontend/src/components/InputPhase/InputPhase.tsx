import { useAppState } from '../../state/useContext';
import { useImageProcessor } from '../../hooks/useImageProcessor';
import { SlackOAuth } from './SlackOAuth';
import { CameraInput } from './CameraInput';

type Props = {
  isSlackLinked: boolean;
  onConnectSlack: () => void;
};

export const InputPhase = ({
  isSlackLinked,
  onConnectSlack,
}: Props) => {
  const { dispatch } = useAppState();

  const handleCapture = (
    blob: Blob,
    shotNumber: string,
  ) => {
    dispatch({
      type: 'SET_PREVIEW',
      file: blob,
      shotNumber,
    });
  };

  const { isCompressing, processImage } = useImageProcessor(
    handleCapture,
    () => {},
  );

  if (!isSlackLinked) {
    return (
      <SlackOAuth
        onConnectSlack={onConnectSlack}
      />
    );
  }

  return (
    <div
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>
        {isCompressing
          ? '画像を処理しています'
          : '画像を撮影'}
      </h2>

      <CameraInput
        onFileSelected={processImage}
        disabled={isCompressing}
      />
    </div>
  );
};
