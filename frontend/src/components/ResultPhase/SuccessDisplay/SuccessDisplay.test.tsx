import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { SuccessDisplay } from './SuccessDisplay';

describe('SuccessDisplay', () => {
  it('renders the success heading', () => {
    render(
      <SuccessDisplay
        onContinue={vi.fn()}
        onExit={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '送信完了',
      }),
    ).toBeInTheDocument();
  });

  it('calls onContinue when the continue button is clicked', () => {
    const onContinue = vi.fn();
    const onExit = vi.fn();

    render(
      <SuccessDisplay
        onContinue={onContinue}
        onExit={onExit}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: 'つづけて撮影',
      }),
    );

    expect(onContinue).toHaveBeenCalledTimes(1);
    expect(onExit).not.toHaveBeenCalled();
  });

  it('calls onExit when the exit button is clicked', () => {
    const onContinue = vi.fn();
    const onExit = vi.fn();

    render(
      <SuccessDisplay
        onContinue={onContinue}
        onExit={onExit}
      />,
    );

    fireEvent.click(
      screen.getByRole('button', {
        name: '終了',
      }),
    );

    expect(onExit).toHaveBeenCalledTimes(1);
    expect(onContinue).not.toHaveBeenCalled();
  });
});
