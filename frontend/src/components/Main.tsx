import { useAppState } from '../state/useContext';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';
import { Spinner } from '../common/Spinner';

export const Main = () => {
  const { state, dispatch } = useAppState();

  switch (state.phase.type) {
    case 'input':
      return (
        <InputPhase
          isSlackLinked={state.isSlackLinked}
          onConnectSlack={() => {
            dispatch({
              type: 'SET_SLACK_LINKED',
              isSlackLinked: true,
            });
          }}
        />
      );

    case 'preview':
      return (
        <PreviewPhase
          blob={state.phase.file}
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
          state={state.phase}
          dispatch={dispatch}
        />
      );

    default: {
      const _exhaustiveCheck: never = state.phase;
      return _exhaustiveCheck;
    }
  }
};
