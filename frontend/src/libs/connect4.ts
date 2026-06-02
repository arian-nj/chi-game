import type { BotDifficulty } from '@/libs/bot-difficulty';
import { pickRandomItem } from '@/libs/bot-difficulty';

export const CONNECT4_ROWS = 6;
export const CONNECT4_COLS = 7;
export const CONNECT4_WIN_LENGTH = 4;

export type Cell = 'R' | 'Y' | null;
export type Player = 'R' | 'Y';

export interface Connect4Settings {
  isBot: boolean;
  botDifficulty: BotDifficulty;
}

export interface Connect4Result {
  winner: Player;
  cells: Array<{ row: number; col: number }>;
}

export type Board = Cell[][];

export function createEmptyBoard(): Board {
  return Array.from({ length: CONNECT4_ROWS }, () =>
    Array.from({ length: CONNECT4_COLS }, () => null),
  );
}

export function cloneBoard(board: Board): Board {
  return board.map(row => [...row]);
}

export function getDropRow(board: Board, col: number): number {
  for (let row = CONNECT4_ROWS - 1; row >= 0; row--) {
    if (board[row]![col] === null) {
      return row;
    }
  }
  return -1;
}

export function getValidColumns(board: Board): number[] {
  return Array.from({ length: CONNECT4_COLS }, (_, col) => col).filter(
    col => getDropRow(board, col) !== -1,
  );
}

export function dropDisc(board: Board, col: number, player: Player): number {
  const row = getDropRow(board, col);
  if (row === -1) {
    return -1;
  }
  board[row]![col] = player;
  return row;
}

export function checkWinnerFrom(board: Board, row: number, col: number): Connect4Result | null {
  const player = board[row]?.[col];
  if (player !== 'R' && player !== 'Y') {
    return null;
  }

  const directions = [
    { rowDelta: 0, colDelta: 1 },
    { rowDelta: 1, colDelta: 0 },
    { rowDelta: 1, colDelta: 1 },
    { rowDelta: 1, colDelta: -1 },
  ];

  for (const { rowDelta, colDelta } of directions) {
    const cells = collectLine(board, row, col, rowDelta, colDelta, player);
    if (cells.length >= CONNECT4_WIN_LENGTH) {
      return { winner: player, cells: cells.slice(0, CONNECT4_WIN_LENGTH) };
    }
  }

  return null;
}

function collectLine(
  board: Board,
  row: number,
  col: number,
  rowDelta: number,
  colDelta: number,
  player: Player,
): Array<{ row: number; col: number }> {
  const cells: Array<{ row: number; col: number }> = [{ row, col }];

  for (const direction of [-1, 1]) {
    let nextRow = row + rowDelta * direction;
    let nextCol = col + colDelta * direction;

    while (
      nextRow >= 0 &&
      nextRow < CONNECT4_ROWS &&
      nextCol >= 0 &&
      nextCol < CONNECT4_COLS &&
      board[nextRow]![nextCol] === player
    ) {
      cells.push({ row: nextRow, col: nextCol });
      nextRow += rowDelta * direction;
      nextCol += colDelta * direction;
    }
  }

  return cells;
}

export function isBoardFull(board: Board): boolean {
  return getValidColumns(board).length === 0;
}

const COLUMN_ORDER = [3, 2, 4, 1, 5, 0, 6];

function orderColumns(columns: number[]): number[] {
  return [...columns].sort(
    (left, right) => COLUMN_ORDER.indexOf(left) - COLUMN_ORDER.indexOf(right),
  );
}

export function getBotMove(
  board: Board,
  bot: Player,
  human: Player,
  difficulty: BotDifficulty = 'hard',
): number {
  const validColumns = getValidColumns(board);
  if (validColumns.length === 0) {
    return -1;
  }

  if (difficulty === 'easy') {
    return pickRandomItem(orderColumns(validColumns));
  }

  const winColumn = findWinningColumn(board, bot);
  if (winColumn !== null) {
    return winColumn;
  }

  const blockColumn = findWinningColumn(board, human);
  if (blockColumn !== null) {
    return blockColumn;
  }

  if (difficulty === 'medium') {
    const preferred = orderColumns(validColumns);
    return pickRandomItem(preferred.slice(0, Math.min(3, preferred.length)));
  }

  let bestColumn = validColumns[0]!;
  let bestScore = Number.NEGATIVE_INFINITY;

  for (const col of orderColumns(validColumns)) {
    const trialBoard = cloneBoard(board);
    const row = dropDisc(trialBoard, col, bot);
    if (row === -1) {
      continue;
    }

    const score = minimax(
      trialBoard,
      6,
      false,
      bot,
      human,
      Number.NEGATIVE_INFINITY,
      Number.POSITIVE_INFINITY,
    );
    if (score > bestScore) {
      bestScore = score;
      bestColumn = col;
    }
  }

  return bestColumn;
}

function findWinningColumn(board: Board, player: Player): number | null {
  for (const col of orderColumns(getValidColumns(board))) {
    const trialBoard = cloneBoard(board);
    const row = dropDisc(trialBoard, col, player);
    if (row !== -1 && checkWinnerFrom(trialBoard, row, col)) {
      return col;
    }
  }
  return null;
}

function minimax(
  board: Board,
  depth: number,
  isMaximizing: boolean,
  bot: Player,
  human: Player,
  alpha: number,
  beta: number,
): number {
  const validColumns = getValidColumns(board);
  if (depth === 0 || validColumns.length === 0) {
    return evaluateBoard(board, bot, human);
  }

  if (isMaximizing) {
    let value = Number.NEGATIVE_INFINITY;
    for (const col of orderColumns(validColumns)) {
      const trialBoard = cloneBoard(board);
      const row = dropDisc(trialBoard, col, bot);
      if (row === -1) {
        continue;
      }

      if (checkWinnerFrom(trialBoard, row, col)) {
        return 1_000_000 - (7 - depth);
      }

      value = Math.max(
        value,
        minimax(trialBoard, depth - 1, false, bot, human, alpha, beta),
      );
      alpha = Math.max(alpha, value);
      if (alpha >= beta) {
        break;
      }
    }
    return value;
  }

  let value = Number.POSITIVE_INFINITY;
  for (const col of orderColumns(validColumns)) {
    const trialBoard = cloneBoard(board);
    const row = dropDisc(trialBoard, col, human);
    if (row === -1) {
      continue;
    }

    if (checkWinnerFrom(trialBoard, row, col)) {
      return -1_000_000 + (7 - depth);
    }

    value = Math.min(value, minimax(trialBoard, depth - 1, true, bot, human, alpha, beta));
    beta = Math.min(beta, value);
    if (alpha >= beta) {
      break;
    }
  }
  return value;
}

function evaluateBoard(board: Board, bot: Player, human: Player): number {
  let score = 0;

  score += evaluateWindows(board, bot, human);
  score += centerColumnBonus(board, bot, human);

  return score;
}

function centerColumnBonus(board: Board, bot: Player, human: Player): number {
  const centerCol = Math.floor(CONNECT4_COLS / 2);
  let botCount = 0;
  let humanCount = 0;

  for (let row = 0; row < CONNECT4_ROWS; row++) {
    if (board[row]![centerCol] === bot) {
      botCount++;
    } else if (board[row]![centerCol] === human) {
      humanCount++;
    }
  }

  return botCount * 6 - humanCount * 6;
}

function evaluateWindows(board: Board, bot: Player, human: Player): number {
  let score = 0;
  const directions = [
    { rowDelta: 0, colDelta: 1 },
    { rowDelta: 1, colDelta: 0 },
    { rowDelta: 1, colDelta: 1 },
    { rowDelta: 1, colDelta: -1 },
  ];

  for (let row = 0; row < CONNECT4_ROWS; row++) {
    for (let col = 0; col < CONNECT4_COLS; col++) {
      for (const { rowDelta, colDelta } of directions) {
        const window = readWindow(board, row, col, rowDelta, colDelta, CONNECT4_WIN_LENGTH);
        if (window.length < CONNECT4_WIN_LENGTH) {
          continue;
        }
        score += scoreWindow(window, bot, human);
      }
    }
  }

  return score;
}

function readWindow(
  board: Board,
  startRow: number,
  startCol: number,
  rowDelta: number,
  colDelta: number,
  length: number,
): Cell[] {
  const window: Cell[] = [];
  for (let index = 0; index < length; index++) {
    const row = startRow + rowDelta * index;
    const col = startCol + colDelta * index;
    if (row < 0 || row >= CONNECT4_ROWS || col < 0 || col >= CONNECT4_COLS) {
      return [];
    }
    window.push(board[row]![col]!);
  }
  return window;
}

function scoreWindow(window: Cell[], bot: Player, human: Player): number {
  const botCount = window.filter(cell => cell === bot).length;
  const humanCount = window.filter(cell => cell === human).length;
  const emptyCount = window.filter(cell => cell === null).length;

  if (botCount > 0 && humanCount > 0) {
    return 0;
  }

  if (botCount === 3 && emptyCount === 1) {
    return 120;
  }
  if (botCount === 2 && emptyCount === 2) {
    return 12;
  }
  if (botCount === 1 && emptyCount === 3) {
    return 2;
  }
  if (humanCount === 3 && emptyCount === 1) {
    return -140;
  }
  if (humanCount === 2 && emptyCount === 2) {
    return -14;
  }
  if (humanCount === 1 && emptyCount === 3) {
    return -3;
  }

  return 0;
}

function isWinningCell(
  winningCells: Array<{ row: number; col: number }>,
  row: number,
  col: number,
): boolean {
  return winningCells.some(cell => cell.row === row && cell.col === col);
}

export { isWinningCell };
