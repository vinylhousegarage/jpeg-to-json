import { Dispatch } from 'react';
import { ResultPhase, AppAction } from '../types';

export const useResultView = (state: ResultPhase, dispatch: Dispatch<AppAction>) => {
  return {
    status: state.status,
    error: state.error,
    onContinue: () => dispatch({ type: 'CONTINUE' }),
    onExit: () => dispatch({ type: 'EXIT' }),
  };
};
