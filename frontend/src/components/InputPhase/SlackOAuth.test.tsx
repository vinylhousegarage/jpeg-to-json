import { render, screen } from '@testing-library/react';
import { fireEvent } from '@testing-library/dom';
import { describe, it, expect, vi } from 'vitest';
import { SlackOAuth } from './SlackOAuth';

describe('SlackOAuth', () => {
  it('should render the unlinked state and call onConnectSlack when the button is clicked', () => {
    const mockOnConnectSlack = vi.fn();

    render(<SlackOAuth isSlackLinked={false} onConnectSlack={mockOnConnectSlack} />);

    // 未連携のメッセージとボタンが表示されているか確認
    expect(screen.getByText(/Slackアカウントを連携します/i)).toBeDefined();
    const button = screen.getByRole('button', { name: /Add to Slack/i });
    expect(button).toBeDefined();

    // ボタンをクリック
    fireEvent.click(button);

    // onConnectSlack が呼ばれたか検証
    expect(mockOnConnectSlack).toHaveBeenCalledTimes(1);
  });

  it('should render the linked state when isSlackLinked is true', () => {
    const mockOnConnectSlack = vi.fn();

    render(<SlackOAuth isSlackLinked={true} onConnectSlack={mockOnConnectSlack} />);

    // 「Slack連携済み」が表示されているか確認
    expect(screen.getByText(/Slack連携済み/i)).toBeDefined();

    // 「Add to Slack」ボタンが表示されていないことを確認
    const button = screen.queryByRole('button', { name: /Add to Slack/i });
    expect(button).toBeNull();
  });
});
