import { useAppState } from '../state/useContext';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';
import { Spinner } from '../common/Spinner';

export const Main = () => {
  const { state, dispatch } = useAppState();

  switch (state.type) {
    case 'input':
      return (
        <InputPhase 
          onRegister={(url: string) => dispatch({ type: 'SUBMIT', slackUrl: url })} 
        />
      );
    case 'preview':
      return (
        <PreviewPhase
          blob={state.file}
          isSending={false}
          onRetake={() => dispatch({ type: 'RETAKE' })}
          onSend={() => dispatch({ type: 'SEND' })}
        />
      );
    case 'upload':
      return <Spinner />;
    case 'result':
      return (
        <ResultPhase 
          state={state} 
          dispatch={dispatch} 
        />
      );
    default:
      {
        const _exhaustiveCheck: never = state;
        return _exhaustiveCheck;
      }
  }
};
