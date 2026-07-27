import { AppPhase, AppAction } from '../types';

export const appReducer = (state: AppPhase, action: AppAction): AppPhase => {
  switch (action.type) {
    case 'SUBMIT':
      return {
        type: 'input',
        slackUrl: action.slackUrl,
      };

    case 'SET_PREVIEW':
      return {
        type: 'preview',
        slackUrl: state.slackUrl ?? '',
        file: action.file,
        previewUrl: action.previewUrl,
        shotNumber: action.shotNumber,
      };

    case 'RETAKE':
      return {
        type: 'input',
        slackUrl: state.slackUrl,
      };

    case 'SEND':
      if (state.type === 'preview') {
        return {
          type: 'upload',
          slackUrl: state.slackUrl,
          shotNumber: state.shotNumber,
        };
      }
      return state;

    case 'EXIT':
      return { type: 'input' };

    case 'START_UPLOAD':
      return {
        type: 'upload',
        slackUrl: state.slackUrl ?? '',
        shotNumber: 'shotNumber' in state ? state.shotNumber : '',
      };

    case 'UPLOAD_COMPLETE':
      return {
        type: 'result',
        slackUrl: state.slackUrl ?? '',
        shotNumber: 'shotNumber' in state ? state.shotNumber : '',
        status: action.status,
        error: action.error
      };

    case 'CONTINUE':
      return {
        type: 'input',
        slackUrl: state.slackUrl,
      };

    default:
      return state;
  }
};
