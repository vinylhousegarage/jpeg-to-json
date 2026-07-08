import { useImageProcessor } from '../../../hooks/useImageProcessor';
import { CameraCaptureProps } from './CameraCapture.types';
import { standardButtonStyle } from '../styles/button';

export const CameraCapture = ({ 
  onCapture, 
  onClearPreview, 
  onSubmit, 
  isSending = false, 
  onError 
}: CameraCaptureProps) => {

  // 圧縮処理中の Loading を管理する状態
  const { isCompressing, processImage } = useImageProcessor(onCapture, onClearPreview, onError);

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    await processImage(file);

    event.target.value = '';
  };

  return (
    <div>
      {/* 「撮り直し」ボタン */}
      <p>
        <label style={standardButtonStyle}>
          撮り直し
          <input
            type="file"
            accept="image/*"
            capture="environment"
            onChange={handleFileChange}
            disabled={isCompressing || isSending}
            style={{ display: 'none' }} // 標準ボタンを隠す
          />
        </label>
      </p>

      {/* 確定・送信ボタン */}
      <p>
        <button
          type="button"
          onClick={onSubmit}
          disabled={isCompressing || isSending}
          style={standardButtonStyle}
        >
          画像を確定し送信
        </button>
      </p>
    </div>
  );
};
