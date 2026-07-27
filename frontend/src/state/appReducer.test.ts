import { appReducer } from './appReducer';
import { AppPhase, AppAction } from '../types';

describe('appReducer', () => {
  const initialState: AppPhase = { type: 'input', slackUrl: 'https://example.com' };

  it('should return the initial state when action is unknown', () => {
    const action = { type: 'INVALID_TYPE' } as unknown as AppAction;
    expect(appReducer(initialState, action)).toEqual(initialState);
  });

  it('should handle SUBMIT', () => {
    const action: AppAction = { type: 'SUBMIT', slackUrl: 'https://new.com' };
    const state = appReducer(initialState, action);
    expect(state).toEqual({ type: 'input', slackUrl: 'https://new.com' });
  });

  it('should handle SET_PREVIEW', () => {
    const file = new File([''], 'test.png');
    const action: AppAction = { 
      type: 'SET_PREVIEW', 
      file,
      previewUrl: 'blob:...',
      shotNumber: 'SHOT-001'
    };
    const state = appReducer(initialState, action);

    if (state.type === 'preview') {
      expect(state.previewUrl).toBe('blob:...');
      expect(state.file).toBe(file);
      expect(state.shotNumber).toBe('SHOT-001');
    } else {
      throw new Error('State should be preview phase');
    }
  });

  it('should handle UPLOAD_COMPLETE', () => {
    const action: AppAction = { 
      type: 'UPLOAD_COMPLETE', 
      status: 'success', 
      error: undefined 
    };
    const previewState: AppPhase = {
      type: 'preview',
      slackUrl: 'https://example.com',
      shotNumber: 'SHOT-001',
      file: new File([''], 'test.png'),
      previewUrl: 'blob:...'
    };
    
    const state = appReducer(previewState, action);
    expect(state).toEqual({
      type: 'result',
      slackUrl: 'https://example.com',
      shotNumber: 'SHOT-001',
      status: 'success',
      error: undefined,
    });
  });

  it('should handle EXIT', () => {
    const action: AppAction = { type: 'EXIT' };
    
    const previewState: AppPhase = { 
      type: 'preview', 
      slackUrl: 'https://example.com',
      shotNumber: 'SHOT-001',
      file: new File([''], 'test.png'),
      previewUrl: 'blob:test'
    };
    
    const state = appReducer(previewState, action);
    
    expect(state).toEqual({ type: 'input' });
  });
});
