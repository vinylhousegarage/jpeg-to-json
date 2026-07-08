import { usePreviewUrl } from '../../../hooks/usePreviewUrl';
import { PreviewAreaProps } from './PreviewArea.types';
import { standardButtonStyle } from '../styles/button';

export const PreviewArea = ({
  blob,
  onRetake,
  onSend,
  isSending
}: PreviewAreaProps) => {
  const imageUrl = usePreviewUrl(blob);

  if (!imageUrl) {
    return null;
  }

  return (
    <div className="preview-container">
      <img
        src={imageUrl}
        alt="preview"
        style={{ maxWidth: '100%', display: 'block', marginBottom: '1rem' }} 
      />

      <div className="button-group" style={{ display: 'flex', gap: '10px' }}>
        {/* 撮り直しボタン */}
        <button
          type="button"
          onClick={onRetake}
          disabled={isSending}
          style={standardButtonStyle}
        >
          撮り直し
        </button>

        {/* 確定し送信ボタン */}
        <button
          type="button"
          onClick={onSend}
          disabled={isSending}
          style={standardButtonStyle}
        >
          画像を確定し送信
        </button>
      </div>
    </div>
  );
};
