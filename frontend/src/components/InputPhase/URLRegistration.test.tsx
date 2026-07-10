import { render, screen } from '@testing-library/react';
import { fireEvent } from '@testing-library/dom';
import { describe, it, expect, vi } from 'vitest';
import { URLRegistration } from './URLRegistration';

describe('URLRegistration', () => {
  it('should call onRegister when a valid URL is submitted', async () => {
    // 関数のモック（実行されたか追跡するための偽関数）を作成
    const mockOnRegister = vi.fn();

    // コンポーネントを描画
    render(<URLRegistration onRegister={mockOnRegister} />);

    // 画面の操作を再現
    const input = screen.getByLabelText(/URLを入力/i);
    const button = screen.getByRole('button', { name: /登録して撮影/i });

    fireEvent.change(input, { target: { value: 'https://hooks.slack.com/test' } });
    fireEvent.click(button);

    // 検証
    expect(mockOnRegister).toHaveBeenCalledWith('https://hooks.slack.com/test');
  });

  it('should display an error message when an invalid URL is submitted', () => {
    render(<URLRegistration onRegister={vi.fn()} />);

    const input = screen.getByLabelText(/URLを入力/i);
    const button = screen.getByRole('button', { name: /登録して撮影/i });

    fireEvent.change(input, { target: { value: 'invalid-url' } });
    fireEvent.click(button);

    // エラーメッセージが表示されているか確認
    expect(screen.getByText(/有効なURL形式で入力してください/i)).toBeDefined();
  });
});
