import React from 'react';
import { standardButtonStyle } from '../../styles/button';

type Props = {
  isSlackLinked: boolean;
  onConnectSlack: () => void;
};

export const SlackOAuth: React.FC<Props> = ({ isSlackLinked, onConnectSlack }) => {
  return (
    <div
      className="slack-oauth"
      style={{
        maxWidth: '375px',
        margin: '0 auto',
        textAlign: 'center'
      }}
    >
      <h2>Slack通知設定</h2>

      {isSlackLinked ? (
        <p style={{ color: 'green', fontWeight: 'bold' }}>
          Slack連携済み
        </p>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px', alignItems: 'center' }}>
          <p style={{ fontSize: '14px', color: '#666' }}>
            Slackアカウントを連携します
          </p>
          <button
            type="button"
            onClick={onConnectSlack}
            style={standardButtonStyle}
          >
            Add to Slack
          </button>
        </div>
      )}
    </div>
  );
};
