import { render, screen } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import { ResultView } from './ResultView';
import { ResultPhase } from '../../types';

describe('ResultView', () => {
  const mockDispatch = vi.fn();

  test('display SuccessDisplay when the status is success', () => {
    const state: ResultPhase = {
      type: 'result',
      slackUrl: 'test-url',
      shotNumber: 'SHOT-001',
      status: 'success'
    };

    render(<ResultView state={state} dispatch={mockDispatch} />);
    
    // 正常形
    expect(screen.getByText('送信完了')).toBeDefined();
  });

  test('display ErrorDisplay when the status is error', () => {
    const state: ResultPhase = {
      type: 'result',
      slackUrl: 'test-url',
      shotNumber: 'SHOT-001',
      status: 'error',
      error: new Error('Network Error')
    };

    render(<ResultView state={state} dispatch={mockDispatch} />);
    
    // 以上系
    expect(screen.getByText('送信失敗')).toBeDefined();
    expect(screen.getByText(/Network Error/)).toBeDefined();
  });
});
