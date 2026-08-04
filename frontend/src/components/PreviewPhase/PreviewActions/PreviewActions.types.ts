export interface PreviewActionsProps {
  onCapture: (
    blob: Blob,
    shotNumber: string,
  ) => void;
  onClearPreview: () => void;
  onSubmit: () => void;
  isSending?: boolean;
  onError?: (error: Error) => void;
}
