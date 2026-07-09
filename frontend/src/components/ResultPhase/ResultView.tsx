import { ResultPhase, AppAction } from '../../types';
import { useResultView } from '../../hooks/useResultView';
import { SuccessDisplay } from './SuccessDisplay';
import { ErrorDisplay } from './ErrorDisplay';

type Props = {
  state: ResultPhase;
  dispatch: React.Dispatch<AppAction>;
};

export const ResultView: React.FC<Props> = ({ state, dispatch }) => {
  const { status, error, onContinue, onExit } = useResultView(state, dispatch);

  return status === 'success' ? (
    <SuccessDisplay onContinue={onContinue} onExit={onExit} />
  ) : (
    <ErrorDisplay error={error} onContinue={onContinue} onExit={onExit} />
  );
};
