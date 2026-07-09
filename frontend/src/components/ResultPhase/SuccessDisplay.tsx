import { standardButtonStyle } from '../PreviewPhase/styles/button';

type Props = {
  onContinue: () => void;
  onExit: () => void;
};

export const SuccessDisplay: React.FC<Props> = ({ onContinue, onExit }) => {
  return (
    <div className="success-display">
      <h2>送信完了</h2>
      <button 
        onClick={onContinue} 
        style={standardButtonStyle}
      >
        つづけて撮影
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
