import React, { useState } from 'react';
import { compressImage } from './utils/compressImage';
import { standardButtonStyle } from './styles/button';

type Props = {
  onCapture: (blob: Blob) => void;
  onClearPreview: () => void; 
  onSubmit: () => void;
  isSending?: boolean;
  onError?: (error: Error) => void;
};

export const CameraCapture: React.FC<Props> = ({ 
  onCapture, 
  onClearPreview, 
  onSubmit, 
  isSending = false, 
  onError 
}) => {
  // 圧縮処理中の Loading を管理する状態
  const [isCompressing, setIsCompressing] = useState(false);

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      // 圧縮処理中の Loading を実行
      setIsCompressing(true);
      
      // プレビューエリアの画像を消去
      onClearPreview();

      // 撮影した画像を圧縮
      const compressedBlob = await compressImage(file);
      
      // 圧縮した画像をプレビューエリアにセット
      onCapture(compressedBlob);
    } catch (error) {
      console.error('Image compression failed:', error);
      if (onError) onError(error instanceof Error ? error : new Error('Failed to compress image'));
    } finally {
      setIsCompressing(false);
      event.target.value = ''; 
    }
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
