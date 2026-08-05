import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { PreviewPhase } from './PreviewPhase';

vi.mock('./PreviewArea', () => ({
  PreviewArea: ({ blob }: { blob: Blob }) => (
    <div data-testid="preview-area">
      {blob.type}
    </div>
  ),
}));

vi.mock('./PreviewActions', () => ({
  PreviewActions: ({
    onRetake,
    onSubmit,
    isSending,
  }: {
    onRetake: () => void;
    onSubmit: () => void;
    isSending?: boolean;
  }) => (
    <div data-testid="preview-actions">
      <button
        type="button"
        onClick={onRetake}
        disabled={isSending}
      >
        撮り直し
      </button>

      <button
        type="button"
        onClick={onSubmit}
        disabled={isSending}
      >
        画像を確定し送信
      </button>
    </div>
  ),
}));

describe('PreviewPhase', () => {
  it('renders the heading, preview area, and preview actions', () => {
    const blob = new Blob(['test-image'], {
      type: 'image/jpeg',
    });

    render(
      <PreviewPhase
        blob={blob}
        onRetake={vi.fn()}
        onSend={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '画像を確認',
      }),
    ).toBeInTheDocument();

    expect(
      screen.getByTestId('preview-area'),
    ).toHaveTextContent('image/jpeg');

    expect(
      screen.getByTestId('preview-actions'),
    ).toBeInTheDocument();
  });

  it('calls onRetake when the retake button is clicked', () => {
    const onRetake = vi.fn();
    const onSend = vi.fn();

    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetake={onRetake}
        onSend={onSend}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    );

    expect(onRetake).toHaveBeenCalledTimes(1);
    expect(onSend).not.toHaveBeenCalled();
  });

  it('calls onSend when the submit button is clicked', () => {
    const onRetake = vi.fn();
    const onSend = vi.fn();

    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetake={onRetake}
        onSend={onSend}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '画像を確定し送信',
      }),
    );

    expect(onSend).toHaveBeenCalledTimes(1);
    expect(onRetake).not.toHaveBeenCalled();
  });

  it('passes isSending to PreviewActions', () => {
    render(
      <PreviewPhase
        blob={new Blob(['test-image'])}
        onRetake={vi.fn()}
        onSend={vi.fn()}
        isSending
      />,
    );

    expect(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    ).toBeDisabled();

    expect(
      screen.getByRole('button', {
        name: '画像を確定し送信',
      }),
    ).toBeDisabled();
  });
});
