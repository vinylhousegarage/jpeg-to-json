import React, { useState } from 'react';
import { uploadUrlSchema } from '../../types/schema';

type Props = {
  onRegister: (url: string) => void;
};

export const URLRegistration: React.FC<Props> = ({ onRegister }) => {
  const [url, setUrl] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Zodによるバリデーション
    const result = uploadUrlSchema.safeParse(url);

    if (!result.success) {
      // エラーメッセージを状態にセット
      setError(result.error.issues[0].message);
      return;
    }

    // バリデーション通過後、親コンポーネントへURLを渡す
    setError(null);
    onRegister(result.data);
  };

  return (
    <div className="url-registration">
      <h2>Slack通知設定</h2>
      <form onSubmit={handleSubmit}>
        <label htmlFor="url">Slack Webhook URL</label>
        <input
          id="url"
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://hooks.slack.com/..."
        />
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <button type="submit">登録して次へ</button>
      </form>
    </div>
  );
};
