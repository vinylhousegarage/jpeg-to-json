import { AppPhase, AppAction } from '../../../types';

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
      };

    case 'RETAKE':
      return {
        type: 'input',
        slackUrl: state.slackUrl,
      };

    case 'SEND':
      if (state.type === 'preview') {
        return {
          type: 'uploading',
          slackUrl: state.slackUrl,
        };
      }
      return state;

    case 'EXIT':
      return { type: 'input' };

    case 'START_UPLOAD':
      return {
        type: 'uploading',
        slackUrl: state.slackUrl ?? '',
      };

    case 'UPLOAD_COMPLETE':
      return {
        type: 'result',
        slackUrl: state.slackUrl ?? '',
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
