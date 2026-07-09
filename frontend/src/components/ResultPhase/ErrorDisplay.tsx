import { standardButtonStyle } from '../../styles/button';

type Props = {
  error?: Error;
  onContinue: () => void;
  onExit: () => void;
};

export const ErrorDisplay: React.FC<Props> = ({ error, onContinue, onExit }) => {
  return (
    <div className="error-display">
      <h2>送信失敗</h2>
      
      {error && <p style={{ color: 'red' }}>{error.message}</p>}
      
      <button
        onClick={onContinue}
        style={standardButtonStyle}
      >
        撮り直し
      </button>
      
      <button
        onClick={onExit}
        style={standardButtonStyle}
      >
        終了
      </button>
    </div>
  );
};
