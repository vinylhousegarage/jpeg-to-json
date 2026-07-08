import { act, renderHook } from '@testing-library/react';
import { vi } from 'vitest';
import { useUrlRegistration } from './useUrlRegistration';

describe('useUrlRegistration', () => {
  it('calls onRegister when handleSubmit is called with a valid URL', () => {
    const mockOnRegister = vi.fn();
    const { result } = renderHook(() => useUrlRegistration(mockOnRegister));

    const mockEvent = {
      preventDefault: vi.fn(),
    } as unknown as React.SubmitEvent<HTMLFormElement>;

    act(() => {
      result.current.setUrl('https://example.com');
    });

    act(() => {
      result.current.handleSubmit(mockEvent);
    });

    expect(mockOnRegister).toHaveBeenCalledWith('https://example.com');
    expect(mockEvent.preventDefault).toHaveBeenCalled();
  });
});
