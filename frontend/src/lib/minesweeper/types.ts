export type Difficulty = 'beginner' | 'intermediate' | 'expert' | 'custom';

export interface GameConfig {
    cellWidth: number;
    cellHeight: number;
    mineCount: number;
}

export interface CustomGameConfig {
    cellWidth: number;
    cellHeight: number;
    mineCount: number;
}

export const DIFFICULTIES: Record<Exclude<Difficulty, 'custom'>, GameConfig> = {
    beginner: { cellWidth: 9, cellHeight: 9, mineCount: 10 },
    intermediate: { cellWidth: 16, cellHeight: 16, mineCount: 40 },
    expert: { cellWidth: 30, cellHeight: 16, mineCount: 99 },
};

export const CUSTOM_LIMITS = {
    minWidth: 9,
    maxWidth: 30,
    minHeight: 9,
    maxHeight: 24,
    minMines: 1,
} as const;

export const DEFAULT_CUSTOM_CONFIG: CustomGameConfig = {
    cellWidth: 16,
    cellHeight: 16,
    mineCount: 40,
};

export function maxMinesForBoard(width: number, height: number): number {
    return Math.max(0, width * height - 9);
}

export function clampCustomConfig(config: CustomGameConfig): CustomGameConfig {
    const cellWidth = Math.min(
        CUSTOM_LIMITS.maxWidth,
        Math.max(CUSTOM_LIMITS.minWidth, Math.round(config.cellWidth)),
    );
    const cellHeight = Math.min(
        CUSTOM_LIMITS.maxHeight,
        Math.max(CUSTOM_LIMITS.minHeight, Math.round(config.cellHeight)),
    );
    const maxMines = maxMinesForBoard(cellWidth, cellHeight);
    const mineCount = Math.min(
        maxMines,
        Math.max(CUSTOM_LIMITS.minMines, Math.round(config.mineCount)),
    );

    return { cellWidth, cellHeight, mineCount };
}
