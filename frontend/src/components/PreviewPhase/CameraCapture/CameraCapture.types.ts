export interface CameraCaptureProps {
  onCapture: (blob: Blob) => void;
  onClearPreview: () => void; 
  onSubmit: () => void;
  isSending?: boolean;
  onError?: (error: Error) => void;
};
