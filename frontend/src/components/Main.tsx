import { useAppState } from '../state/useContext';
import { usePresignUpload } from '../hooks/usePresignUpload';
import { InputPhase } from './InputPhase';
import { PreviewPhase } from './PreviewPhase';
import { ResultPhase } from './ResultPhase';
import { Spinner } from '../common/Spinner';

export const Main = () => {
  const { state, dispatch } = useAppState();
  const { send } = usePresignUpload(dispatch);

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

    case 'preview': {
      const { file, shotNumber } = state.phase;

      return (
        <PreviewPhase
          blob={file}
          isSending={false}
          onRetake={() => dispatch({ type: 'RETAKE' })}
          onSend={() =>
            send(
              file,
              shotNumber,
            )
          }
        />
      );
    }

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
