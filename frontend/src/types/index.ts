export type InputPhase = {
  type: 'input';
  isLinked: boolean;
};

export type PreviewPhase = {
  type: 'preview';
  shotNumber: string;
  file: Blob;
  previewUrl: string;
};

export type UploadPhase = {
  type: 'upload';
  shotNumber: string;
};

export type ResultPhase = {
  type: 'result';
  shotNumber: string;
  status: 'success' | 'error';
  error?: Error;
};

export type AppPhase = InputPhase | PreviewPhase | UploadPhase | ResultPhase;

  export type AppAction =
  | { type: 'SET_LINKED'; isLinked: boolean }
  | { type: 'SET_PREVIEW'; file: Blob; previewUrl: string; shotNumber: string }
  | { type: 'RETAKE' }
  | { type: 'SEND' }
  | { type: 'EXIT' }
  | { type: 'START_UPLOAD' }
  | { type: 'UPLOAD_COMPLETE'; status: 'success' | 'error'; error?: Error }
  | { type: 'CONTINUE' };
