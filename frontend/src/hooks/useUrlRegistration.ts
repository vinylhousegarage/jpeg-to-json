import { useState } from 'react';
import { uploadUrlSchema } from '../types/schema';

export const useUrlRegistration = (onRegister: (url: string) => void) => {
  const [url, setUrl] = useState('');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    const result = uploadUrlSchema.safeParse(url);
    if (!result.success) {
      setError(result.error.issues[0].message);
      return;
    }
    setError(null);
    onRegister(result.data);
  };

  return { url, setUrl, error, handleSubmit };
};
