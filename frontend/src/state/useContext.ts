import { useContext } from 'react';
import { AppContext } from '../state/AppContext';
import { AppPhase, AppAction } from '../types';
import { Dispatch } from 'react';

type AppContextType = {
  state: AppPhase;
  dispatch: Dispatch<AppAction>;
};

export const useAppState = () => {
  const context = useContext(AppContext) as AppContextType | undefined;
  if (context === undefined) {
    throw new Error('useAppState must be used within an AppProvider');
  }
  return context;
};
