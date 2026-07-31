import { AppPhase, AppAction } from '../types';

export const appReducer = (state: AppPhase, action: AppAction): AppPhase => {
  switch (action.type) {
    case 'SET_LINKED':
      return {
        type: 'input',
        isLinked: action.isLinked,
      };

    case 'SET_PREVIEW':
      return {
        type: 'preview',
        file: action.file,
        previewUrl: action.previewUrl,
        shotNumber: action.shotNumber,
      };

    case 'RETAKE':
    case 'CONTINUE':
    case 'EXIT':
      return {
        type: 'input',
        isLinked: 'isLinked' in state ? state.isLinked : false,
      };

    case 'SEND':
      if (state.type === 'preview') {
        return {
          type: 'upload',
          shotNumber: state.shotNumber,
        };
      }
      return state;

    case 'START_UPLOAD':
      return {
        type: 'upload',
        shotNumber: 'shotNumber' in state ? state.shotNumber : '',
      };

    case 'UPLOAD_COMPLETE':
      return {
        type: 'result',
        shotNumber: 'shotNumber' in state ? state.shotNumber : '',
        status: action.status,
        error: action.error,
      };

    default:
      return state;
  }
};
