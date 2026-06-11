import { BOX_SIZE, GRID_SIZE } from './types';

export type CellValue = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9;
export type Grid = CellValue[][];

function shuffle<T>(items: T[]): T[] {
    const arr = [...items];
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j]!, arr[i]!];
    }
    return arr;
}

export function createEmptyGrid(): Grid {
    return Array.from({ length: GRID_SIZE }, () => Array(GRID_SIZE).fill(0) as CellValue[]);
}

export function isValidPlacement(grid: Grid, row: number, col: number, value: CellValue): boolean {
    if (value === 0) return true;

    for (let c = 0; c < GRID_SIZE; c++) {
        if (c !== col && grid[row]![c] === value) return false;
    }
    for (let r = 0; r < GRID_SIZE; r++) {
        if (r !== row && grid[r]![col] === value) return false;
    }

    const boxRow = Math.floor(row / BOX_SIZE) * BOX_SIZE;
    const boxCol = Math.floor(col / BOX_SIZE) * BOX_SIZE;
    for (let r = 0; r < BOX_SIZE; r++) {
        for (let c = 0; c < BOX_SIZE; c++) {
            const nr = boxRow + r;
            const nc = boxCol + c;
            if ((nr !== row || nc !== col) && grid[nr]![nc] === value) return false;
        }
    }
    return true;
}

function solveGrid(grid: Grid): boolean {
    for (let row = 0; row < GRID_SIZE; row++) {
        for (let col = 0; col < GRID_SIZE; col++) {
            if (grid[row]![col] !== 0) continue;

            for (const num of shuffle([1, 2, 3, 4, 5, 6, 7, 8, 9] as CellValue[])) {
                if (!isValidPlacement(grid, row, col, num)) continue;
                grid[row]![col] = num;
                if (solveGrid(grid)) return true;
                grid[row]![col] = 0;
            }
            return false;
        }
    }
    return true;
}

export function generateSolution(): Grid {
    const grid = createEmptyGrid();
    solveGrid(grid);
    return grid;
}

export function createPuzzle(solution: Grid, givenCount: number): Grid {
    const puzzle = solution.map(row => [...row]) as Grid;
    const positions = shuffle(
        Array.from({ length: GRID_SIZE * GRID_SIZE }, (_, i) => [
            Math.floor(i / GRID_SIZE),
            i % GRID_SIZE,
        ] as [number, number]),
    );

    const cellsToRemove = GRID_SIZE * GRID_SIZE - givenCount;
    for (let i = 0; i < cellsToRemove && i < positions.length; i++) {
        const [row, col] = positions[i]!;
        puzzle[row]![col] = 0;
    }
    return puzzle;
}

export function generatePuzzle(givenCount: number): { puzzle: Grid; solution: Grid } {
    const solution = generateSolution();
    const puzzle = createPuzzle(solution, givenCount);
    return { puzzle, solution };
}

export function isGridComplete(grid: Grid): boolean {
    for (let row = 0; row < GRID_SIZE; row++) {
        for (let col = 0; col < GRID_SIZE; col++) {
            const value = grid[row]![col]!;
            if (value === 0 || !isValidPlacement(grid, row, col, value)) return false;
        }
    }
    return true;
}

export function getConflictingCells(grid: Grid): Set<string> {
    const conflicts = new Set<string>();

    function markDuplicates(getCells: () => [number, number][]) {
        const cells = getCells();
        const byValue = new Map<CellValue, [number, number][]>();
        for (const [row, col] of cells) {
            const value = grid[row]![col]!;
            if (value === 0) continue;
            const list = byValue.get(value) ?? [];
            list.push([row, col]);
            byValue.set(value, list);
        }
        for (const positions of byValue.values()) {
            if (positions.length > 1) {
                for (const [row, col] of positions) {
                    conflicts.add(`${row},${col}`);
                }
            }
        }
    }

    for (let row = 0; row < GRID_SIZE; row++) {
        markDuplicates(() => Array.from({ length: GRID_SIZE }, (_, col) => [row, col] as [number, number]));
    }
    for (let col = 0; col < GRID_SIZE; col++) {
        markDuplicates(() => Array.from({ length: GRID_SIZE }, (_, row) => [row, col] as [number, number]));
    }
    for (let boxRow = 0; boxRow < GRID_SIZE; boxRow += BOX_SIZE) {
        for (let boxCol = 0; boxCol < GRID_SIZE; boxCol += BOX_SIZE) {
            markDuplicates(() => {
                const cells: [number, number][] = [];
                for (let r = 0; r < BOX_SIZE; r++) {
                    for (let c = 0; c < BOX_SIZE; c++) {
                        cells.push([boxRow + r, boxCol + c]);
                    }
                }
                return cells;
            });
        }
    }

    return conflicts;
}
