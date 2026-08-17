import { standardButtonStyle } from '../../../styles/button';

type Props = {
  onRetake: () => void;
  onSubmit: () => void;
  isSending?: boolean;
};

export const PreviewActions = ({
  onRetake,
  onSubmit,
  isSending = false,
}: Props) => {
  return (
    <div
      className="button-group"
      style={{
        display: 'flex',
        justifyContent: 'center',
        gap: '10px',
        width: '100%',
      }}
    >
      <button
        type="button"
        onClick={onRetake}
        disabled={isSending}
        style={standardButtonStyle}
      >
        撮り直し
      </button>

      <button
        type="button"
        onClick={onSubmit}
        disabled={isSending}
        style={standardButtonStyle}
      >
        画像を確定し送信
      </button>
    </div>
  );
};
