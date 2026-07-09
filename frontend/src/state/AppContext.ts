import { createContext } from 'react';
import { AppPhase, AppAction } from '../types';

export type AppContextType = {
  state: AppPhase;
  dispatch: React.Dispatch<AppAction>;
};

// 初期値として undefined を許容しフック内でチェック
export const AppContext = createContext<AppContextType | undefined>(undefined);
