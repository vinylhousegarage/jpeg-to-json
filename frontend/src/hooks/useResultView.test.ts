import { describe, test, expect, vi } from 'vitest';
import { useResultView } from './useResultView';
import { ResultPhase, AppAction } from '../types';
import { Dispatch } from 'react';

describe('useResultView', () => {
  test('onContinue should dispatch CONTINUE action', () => {
    const mockDispatch = vi.fn() as unknown as Dispatch<AppAction>;
    
    const state: ResultPhase = { 
      type: 'result', 
      slackUrl: 'https://hooks.slack.com/test',
      status: 'success'
    };
    
    const hook = useResultView(state, mockDispatch);
    hook.onContinue();

    expect(mockDispatch).toHaveBeenCalledWith({ type: 'CONTINUE' });
  });

  test('onExit should dispatch EXIT action', () => {
    const mockDispatch = vi.fn() as unknown as Dispatch<AppAction>;
    
    const state: ResultPhase = { 
      type: 'result', 
      slackUrl: 'https://hooks.slack.com/test',
      status: 'error',
      error: new Error('Upload failed')
    };
    
    const hook = useResultView(state, mockDispatch);
    hook.onExit();

    expect(mockDispatch).toHaveBeenCalledWith({ type: 'EXIT' });
  });
});
