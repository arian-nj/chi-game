export type BotDifficulty = 'easy' | 'medium' | 'hard';

export const botDifficulties: BotDifficulty[] = ['easy', 'medium', 'hard'];

export function pickRandomItem<T>(items: T[]): T {
  return items[Math.floor(Math.random() * items.length)]!;
}
