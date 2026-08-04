import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { InputPhase } from './InputPhase';
import * as UseContextModule from '../../state/useContext';
import * as UseImageProcessorModule from '../../hooks/useImageProcessor';
import * as CreateShotNumberModule from '../../utils/createShotNumber';

const mockDispatch = vi.fn();
const mockProcessImage = vi.fn();

vi.mock('../../state/useContext', () => ({
  useAppState: vi.fn(),
}));

vi.mock('../../hooks/useImageProcessor', () => ({
  useImageProcessor: vi.fn(),
}));

vi.mock('../../utils/createShotNumber', () => ({
  createShotNumber: vi.fn(),
}));

vi.mock('./SlackOAuth', () => ({
  SlackOAuth: ({
    onConnectSlack,
  }: {
    onConnectSlack: () => void;
  }) => (
    <button type="button" onClick={onConnectSlack}>
      DM通知を設定
    </button>
  ),
}));

vi.mock('./CameraInput', () => ({
  CameraInput: ({
    disabled,
  }: {
    onFileSelected: (file: File) => Promise<void>;
    disabled?: boolean;
  }) => (
    <input
      aria-label="カメラを起動"
      type="file"
      disabled={disabled}
    />
  ),
}));

describe('InputPhase', () => {
  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(
      UseContextModule.useAppState,
    ).mockReturnValue({
      state: {
        isSlackLinked: false,
        phase: {
          type: 'input',
        },
      },
      dispatch: mockDispatch,
    });

    vi.mocked(
      UseImageProcessorModule.useImageProcessor,
    ).mockReturnValue({
      isCompressing: false,
      processImage: mockProcessImage,
    });

    vi.mocked(
      CreateShotNumberModule.createShotNumber,
    ).mockReturnValue('SHOT-001');
  });

  it('renders SlackOAuth when Slack is not linked', () => {
    render(
      <InputPhase
        isSlackLinked={false}
        onConnectSlack={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('button', {
        name: 'DM通知を設定',
      }),
    ).toBeInTheDocument();

    expect(
      screen.queryByLabelText('カメラを起動'),
    ).not.toBeInTheDocument();
  });

  it('renders CameraInput when Slack is linked', () => {
    render(
      <InputPhase
        isSlackLinked
        onConnectSlack={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '画像を撮影',
      }),
    ).toBeInTheDocument();

    expect(
      screen.getByLabelText('カメラを起動'),
    ).toBeEnabled();
  });

  it('shows processing state and disables CameraInput', () => {
    vi.mocked(
      UseImageProcessorModule.useImageProcessor,
    ).mockReturnValue({
      isCompressing: true,
      processImage: mockProcessImage,
    });

    render(
      <InputPhase
        isSlackLinked
        onConnectSlack={vi.fn()}
      />,
    );

    expect(
      screen.getByRole('heading', {
        name: '画像を処理しています',
      }),
    ).toBeInTheDocument();

    expect(
      screen.getByLabelText('カメラを起動'),
    ).toBeDisabled();
  });

  it('dispatches SET_PREVIEW when capture completes', () => {
    const blob = new Blob(['image'], {
      type: 'image/jpeg',
    });

    vi.mocked(
      UseImageProcessorModule.useImageProcessor,
    ).mockImplementation((onCapture) => {
      onCapture(blob);

      return {
        isCompressing: false,
        processImage: mockProcessImage,
      };
    });

    render(
      <InputPhase
        isSlackLinked
        onConnectSlack={vi.fn()}
      />,
    );

    expect(
      CreateShotNumberModule.createShotNumber,
    ).toHaveBeenCalledTimes(1);

    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'SET_PREVIEW',
      file: blob,
      shotNumber: 'SHOT-001',
    });
  });
});
