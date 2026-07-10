import { useUrlRegistration } from '../../hooks/useUrlRegistration';

type Props = {
  onRegister: (url: string) => void;
};

export const URLRegistration: React.FC<Props> = ({ onRegister }) => {
  const { url, setUrl, error, handleSubmit } = useUrlRegistration(onRegister);

  return (
    <div 
      className="url-registration" 
      style={{ 
        maxWidth: '375px', 
        margin: '0 auto', 
        textAlign: 'center'
      }}
    >
      <h2>Slack通知設定</h2>
      <form 
        onSubmit={handleSubmit} 
        style={{ 
          display: 'flex', 
          flexDirection: 'column', 
          gap: '10px',
          alignItems: 'center'
        }}
      >
        <label htmlFor="url">URLを入力</label>
        <input
          id="url"
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://hooks.slack.com/..."
          style={{ width: '100%', boxSizing: 'border-box' }}
        />
        {error && <p style={{ color: 'red', margin: '0' }}>{error}</p>}
        <button type="submit">登録して撮影</button>
      </form>
    </div>
  );
};
