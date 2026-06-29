import { expect, it } from 'vitest';
import { add } from './math';

it('1 + 1 は 2 になること', () => {
  expect(add(1, 1)).toBe(2);
});
