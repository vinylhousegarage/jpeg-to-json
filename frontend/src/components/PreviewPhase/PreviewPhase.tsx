import { PreviewArea } from './PreviewArea';
import { PreviewActions } from './PreviewActions';

type Props = {
  blob: Blob;
  onRetake: () => void;
  onSend: () => void;
  isSending?: boolean;
};

export const PreviewPhase = ({
  blob,
  onRetake,
  onSend,
  isSending = false,
}: Props) => {
  return (
    <div
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center',
      }}
    >
      <h2>画像を確認</h2>

      <PreviewArea blob={blob} />

      <PreviewActions
        onRetake={onRetake}
        onSubmit={onSend}
        isSending={isSending}
      />
    </div>
  );
};
