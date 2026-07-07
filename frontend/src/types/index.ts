export type InputPhase = { type: 'input' };

export type PreviewPhase = {
  type: 'preview';
  file: Blob;
  previewUrl: string;
};

export type UploadingPhase = {
   type:   'uploading';
   status: UploadState;
   error?: Error;
};

export type AppPhase = InputPhase | PreviewPhase | UploadingPhase;

export type UploadState = 'idle' | 'uploading' | 'success' | 'error';

export type AppAction =
  | { type: 'SUBMIT' }
  | { type: 'SET_PREVIEW'; file: Blob; previewUrl: string }
  | { type: 'RETAKE' }
  | { type: 'SEND' }
  | { type: 'EXIT' }
  | { type: 'SET_UPLOAD_STATUS'; status: UploadState }
  | { type: 'CONTINUE' };
