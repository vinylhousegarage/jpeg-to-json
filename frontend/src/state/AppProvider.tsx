import { useReducer, ReactNode } from 'react';
import { AppContext } from './AppContext';
import { appReducer } from './appReducer';
import { AppPhase } from '../types';

const initialState: AppPhase = { type: 'input', isLinked: false };

export const AppProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(appReducer, initialState);

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
};
