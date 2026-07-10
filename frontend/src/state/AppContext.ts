import { createContext, useContext } from 'react';
import { AppPhase, AppAction } from '../types';

export type AppContextType = {
  state: AppPhase;
  dispatch: React.Dispatch<AppAction>;
};

// 初期値として undefined を許容しフック内でチェック
export const AppContext = createContext<AppContextType | undefined>(undefined);

export const useAppState = () => {
  const context = useContext(AppContext);
  
  if (context === undefined) {
    throw new Error('useAppState must be used within an AppProvider');
  }
  
  return context;
};
