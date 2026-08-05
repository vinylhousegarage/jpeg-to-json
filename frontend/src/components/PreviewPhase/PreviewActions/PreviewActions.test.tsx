import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { PreviewActions } from './PreviewActions';

describe('PreviewActions', () => {
  it('calls onRetake when the retake button is clicked', () => {
    const onRetake = vi.fn();
    const onSubmit = vi.fn();

    render(
      <PreviewActions
        onRetake={onRetake}
        onSubmit={onSubmit}
      />,
    );

    const retakeButton = screen.getByRole('button', {
      name: '撮り直し',
    });

    fireEvent.click(retakeButton);

    expect(onRetake).toHaveBeenCalledTimes(1);
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('calls onSubmit when the submit button is clicked', () => {
    const onRetake = vi.fn();
    const onSubmit = vi.fn();

    render(
      <PreviewActions
        onRetake={onRetake}
        onSubmit={onSubmit}
      />,
    );

    const submitButton = screen.getByRole('button', {
      name: '画像を確定し送信',
    });

    fireEvent.click(submitButton);

    expect(onSubmit).toHaveBeenCalledTimes(1);
    expect(onRetake).not.toHaveBeenCalled();
  });

  it('disables both buttons when isSending is true', () => {
    const onRetake = vi.fn();
    const onSubmit = vi.fn();

    render(
      <PreviewActions
        onRetake={onRetake}
        onSubmit={onSubmit}
        isSending
      />,
    );

    const retakeButton = screen.getByRole('button', {
      name: '撮り直し',
    });

    const submitButton = screen.getByRole('button', {
      name: '画像を確定し送信',
    });

    expect(retakeButton).toBeDisabled();
    expect(submitButton).toBeDisabled();

    fireEvent.click(retakeButton);
    fireEvent.click(submitButton);

    expect(onRetake).not.toHaveBeenCalled();
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('enables both buttons by default', () => {
    render(
      <PreviewActions
        onRetake={vi.fn()}
        onSubmit={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('button', {
        name: '撮り直し',
      }),
    ).toBeEnabled();

    expect(
      screen.getByRole('button', {
        name: '画像を確定し送信',
      }),
    ).toBeEnabled();
  });
});
