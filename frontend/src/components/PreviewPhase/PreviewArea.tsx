import { useState, useEffect } from 'react';
import { standardButtonStyle } from './styles/button';
import { PreviewAreaProps } from './PreviewArea.types';

export const PreviewArea = ({
  blob,
  onRetake,
  onSend,
  isSending
}: PreviewAreaProps) => {
  const [imageUrl, setImageUrl] = useState<string>('');

  useEffect(() => {
    const url = URL.createObjectURL(blob);
    setImageUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [blob]);

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
