export type InputPhase = { type: 'input' };

export type PreviewPhase = {
  type: 'preview';
  file: Blob;
  previewUrl: string;
};

export type UploadingPhase = { type: 'uploading' };

export type AppPhase = InputPhase | PreviewPhase | UploadingPhase;

export type AppAction =
  | { type: 'SET_PREVIEW'; file: Blob; previewUrl: string }
  | { type: 'START_UPLOAD' }
  | { type: 'RESET' };
