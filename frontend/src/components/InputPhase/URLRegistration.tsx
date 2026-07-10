import { useUrlRegistration } from '../../hooks/useUrlRegistration';

type Props = {
  onRegister: (url: string) => void;
};

export const URLRegistration: React.FC<Props> = ({ onRegister }) => {
  const { url, setUrl, error, handleSubmit } = useUrlRegistration(onRegister);

  return (
    <div className="url-registration">
      <h2>Slack通知設定</h2>
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
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
