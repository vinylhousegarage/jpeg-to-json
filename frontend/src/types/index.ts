export type InputPhase = {
  type: 'input';
  slackUrl?: string;
};

export type PreviewPhase = {
  type: 'preview';
  slackUrl: string;
  file: Blob;
  previewUrl: string;
};

export type UploadingPhase = {
  type: 'uploading';
  slackUrl: string;
};

export type ResultPhase = {
  type: 'result';
  slackUrl: string;
  status: 'success' | 'error';
  error?: Error;
}

export type AppPhase = InputPhase | PreviewPhase | UploadingPhase | ResultPhase;

export type AppAction =
  | { type: 'SUBMIT'; slackUrl: string }
  | { type: 'SET_PREVIEW'; file: Blob; previewUrl: string }
  | { type: 'RETAKE' }
  | { type: 'SEND' }
  | { type: 'EXIT' }
  | { type: 'START_UPLOAD' }
  | { type: 'UPLOAD_COMPLETE'; status: 'success' | 'error'; error?: Error }
  | { type: 'CONTINUE' };
