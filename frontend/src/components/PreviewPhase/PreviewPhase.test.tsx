import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { PreviewPhase } from './PreviewPhase';
import type { CameraCaptureProps } from './CameraCapture';

// テスト用の拡張型を定義
interface MockCameraCaptureProps extends CameraCaptureProps {
  onRetake?: () => void;
  onExit?: () => void;
}

// モック定義
vi.mock('./CameraCapture', () => ({
  // 拡張した型を適用
  CameraCapture: (props: MockCameraCaptureProps) => (
    <div>
      <button 
        data-testid="mock-retake" 
        onClick={() => props.onRetake?.()}
      >
        撮り直し
      </button>
      
      <button 
        data-testid="mock-submit" 
        onClick={props.onSubmit} 
        disabled={props.isSending}
      >
        画像を確定し送信
      </button>

      <button 
        data-testid="mock-exit" 
        onClick={() => props.onExit?.()}
      >
        終了
      </button>
    </div>
  ),
}));

vi.mock('./PreviewArea', () => ({
  PreviewArea: () => <div data-testid="mock-preview" />,
}));

describe('PreviewPhase', () => {

  describe('UI & Interaction', () => {
    // 正常系：渡された画像がプレビューに表示されるか
    it('should render the image provided in props', () => {
      const mockBlob = new Blob();
      render(<PreviewPhase initialBlob={mockBlob} />);
      expect(screen.getByTestId('mock-preview')).toBeInTheDocument();
    });

    // 正常系：送信処理が正しく動作するか
    it('should handle submission loading and reset', async () => {
      render(<PreviewPhase />);
      fireEvent.click(screen.getByTestId('mock-retake')); // 撮り直しボタンをクリック

      const submitBtn = screen.getByTestId('mock-submit');

      // 送信開始
      fireEvent.click(submitBtn);

      // ローディング中の状態を確認 (送信ボタンがクリックされた瞬間)
      expect(submitBtn).toBeDisabled();

      // 送信完了後の状態を確認 (非同期完了を待機)
      await waitFor(() => {
        expect(submitBtn).not.toBeDisabled();
      });
    });
  });

  // URL.revokeObjectURL をモック化
  const revokeObjectURLSpy = vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {});

  describe('Memory Management', () => {
    beforeEach(() => {
      vi.clearAllMocks();
    });

    afterAll(() => {
      revokeObjectURLSpy.mockRestore();
    });

    // 撮り直し時のテスト：関数呼び出しとメモリ解放
    it('should call onRetake and revoke blob URL when retake button is clicked', () => {
      const mockOnRetake = vi.fn();
      const mockBlob = new Blob();
      const blobUrl = 'blob:test-url';

      render(<PreviewPhase initialBlob={mockBlob} blobUrl={blobUrl} onRetake={mockOnRetake} />);

      fireEvent.click(screen.getByTestId('mock-retake'));

      // 適切な関数が呼ばれたか
      expect(mockOnRetake).toHaveBeenCalled();
      // メモリ解放が呼ばれたか
      expect(revokeObjectURLSpy).toHaveBeenCalledWith(blobUrl);
    });

    // 送信完了時のテスト：API完了後のメモリ解放
    it('should revoke blob URL after successful submission', async () => {
      const mockOnSubmit = vi.fn().mockResolvedValue({ status: 200 });
      const blobUrl = 'blob:test-url';

      render(<PreviewPhase initialBlob={new Blob()} blobUrl={blobUrl} onSubmit={mockOnSubmit} />);

      // 送信ボタンをクリック
      fireEvent.click(screen.getByTestId('mock-submit'));

      // 送信処理（非同期）が終わるまで待機
      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
      });

      // 送信成功後にメモリ解放が呼ばれたか
      expect(revokeObjectURLSpy).toHaveBeenCalledWith(blobUrl);
    });

    it('should restore button and revoke blob URL even when submission fails', async () => {
      // エラーを返すモック
      const mockOnSubmit = vi.fn().mockRejectedValue(new Error('Network Error'));
      const blobUrl = 'blob:test-url';

      render(<PreviewPhase initialBlob={new Blob()} blobUrl={blobUrl} onSubmit={mockOnSubmit} />);

      const submitBtn = screen.getByTestId('mock-submit');

      // 送信ボタンをクリック
      fireEvent.click(submitBtn);

      // 非同期処理の完了を待機
      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
      });

      // UIが「送信中」状態から戻っているか（ボタンが押せるようになっているか）
      expect(submitBtn).not.toBeDisabled();

      // エラー発生時にメモリが解放されたか
      expect(revokeObjectURLSpy).toHaveBeenCalledWith(blobUrl);
    });
  });
});
