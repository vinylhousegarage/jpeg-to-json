import { useAppState } from '../state/AppContext';
import { SlackOAuth } from './InputPhase/SlackOAuth';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';
import { Spinner } from '../common/Spinner';

export const Main = () => {
  const { state, dispatch } = useAppState();

  switch (state.type) {
    case 'input':
      return (
        <SlackOAuth 
          isSlackLinked={state.isSlackLinked} 
          onConnectSlack={() => {
            dispatch({ type: 'SET_SLACK_LINKED', isSlackLinked: true });
          }} 
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
