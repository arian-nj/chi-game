export type Difficulty = 'easy' | 'medium' | 'hard';

export const DIFFICULTIES: Record<Difficulty, { givenCount: number; label: string }> = {
    easy: { givenCount: 40, label: 'Easy' },
    medium: { givenCount: 32, label: 'Medium' },
    hard: { givenCount: 24, label: 'Hard' },
};

export const GRID_SIZE = 9;
export const BOX_SIZE = 3;
