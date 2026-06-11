export type Difficulty = 'beginner' | 'intermediate' | 'expert';

export const DIFFICULTIES: Record<Difficulty, { cellWidth: number; cellHeight: number; mineCount: number }> = {
    beginner: { cellWidth: 9, cellHeight: 9, mineCount: 10 },
    intermediate: { cellWidth: 16, cellHeight: 16, mineCount: 40 },
    expert: { cellWidth: 30, cellHeight: 16, mineCount: 99 },
};
