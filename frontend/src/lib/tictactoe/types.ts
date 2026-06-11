export type BoardSize = 3 | 5;
export type BotDifficulty = 'easy' | 'medium' | 'hard';
export type GameMode = 'local' | 'bot';

export const BOARD_SIZES: Record<'classic' | 'expert', BoardSize> = {
    classic: 3,
    expert: 5,
};

export type BoardSizeKey = keyof typeof BOARD_SIZES;
