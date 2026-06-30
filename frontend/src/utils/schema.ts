import { z } from 'zod';

export const uploadUrlSchema = z.string()
  .min(1, { message: "URLを入力してください" })
  .refine((val) => {
    try {
      new URL(val);
      return true;
    } catch {
      return false;
    }
  }, { message: "有効なURL形式で入力してください" });
