export type Cell = 'X' | 'O' | null;
export type Player = 'X' | 'O';
export type BoardSize = 3 | 5;

export interface TicTacToeSettings {
  isBot: boolean;
  boardSize: BoardSize;
}

export interface GameResult {
  winner: Player;
  line: number[];
}

export function winLengthForSize(size: BoardSize): number {
  return size === 3 ? 3 : 4;
}

export function createEmptyBoard(size: BoardSize): Cell[] {
  return Array.from({ length: size * size }, () => null);
}

export function getAvailableMoves(board: Cell[]): number[] {
  return board.flatMap((cell, index) => (cell === null ? [index] : []));
}

export function checkWinner(board: Cell[], size: BoardSize): GameResult | null {
  const winLength = winLengthForSize(size);

  for (let row = 0; row < size; row++) {
    for (let col = 0; col <= size - winLength; col++) {
      const line = Array.from({ length: winLength }, (_, offset) => row * size + col + offset);
      const winner = lineWinner(board, line);
      if (winner) {
        return { winner, line };
      }
    }
  }

  for (let col = 0; col < size; col++) {
    for (let row = 0; row <= size - winLength; row++) {
      const line = Array.from({ length: winLength }, (_, offset) => (row + offset) * size + col);
      const winner = lineWinner(board, line);
      if (winner) {
        return { winner, line };
      }
    }
  }

  for (let row = 0; row <= size - winLength; row++) {
    for (let col = 0; col <= size - winLength; col++) {
      const mainDiagonal = Array.from(
        { length: winLength },
        (_, offset) => (row + offset) * size + (col + offset),
      );
      const mainWinner = lineWinner(board, mainDiagonal);
      if (mainWinner) {
        return { winner: mainWinner, line: mainDiagonal };
      }

      const antiDiagonal = Array.from(
        { length: winLength },
        (_, offset) => (row + offset) * size + (col + winLength - 1 - offset),
      );
      const antiWinner = lineWinner(board, antiDiagonal);
      if (antiWinner) {
        return { winner: antiWinner, line: antiDiagonal };
      }
    }
  }

  return null;
}

function lineWinner(board: Cell[], line: number[]): Player | null {
  const first = board[line[0]!];
  if (first !== 'X' && first !== 'O') {
    return null;
  }

  return line.every(index => board[index] === first) ? first : null;
}

export function isDraw(board: Cell[], size: BoardSize): boolean {
  return checkWinner(board, size) === null && getAvailableMoves(board).length === 0;
}

export function getBotMove(board: Cell[], size: BoardSize, bot: Player, human: Player): number {
  const moves = getAvailableMoves(board);
  if (moves.length === 0) {
    return -1;
  }

  if (size === 3) {
    return getPerfectMove(board, size, bot, human);
  }

  return getHeuristicMove(board, size, bot, human, moves);
}

function getPerfectMove(board: Cell[], size: BoardSize, bot: Player, human: Player): number {
  let bestScore = Number.NEGATIVE_INFINITY;
  let bestMove = getAvailableMoves(board)[0]!;

  for (const move of getAvailableMoves(board)) {
    board[move] = bot;
    const score = minimax(board, size, 0, false, bot, human);
    board[move] = null;

    if (score > bestScore) {
      bestScore = score;
      bestMove = move;
    }
  }

  return bestMove;
}

function minimax(
  board: Cell[],
  size: BoardSize,
  depth: number,
  isMaximizing: boolean,
  bot: Player,
  human: Player,
): number {
  const result = checkWinner(board, size);
  if (result?.winner === bot) {
    return 10 - depth;
  }
  if (result?.winner === human) {
    return depth - 10;
  }
  if (isDraw(board, size)) {
    return 0;
  }

  if (isMaximizing) {
    let bestScore = Number.NEGATIVE_INFINITY;
    for (const move of getAvailableMoves(board)) {
      board[move] = bot;
      bestScore = Math.max(bestScore, minimax(board, size, depth + 1, false, bot, human));
      board[move] = null;
    }
    return bestScore;
  }

  let bestScore = Number.POSITIVE_INFINITY;
  for (const move of getAvailableMoves(board)) {
    board[move] = human;
    bestScore = Math.min(bestScore, minimax(board, size, depth + 1, true, bot, human));
    board[move] = null;
  }
  return bestScore;
}

function getHeuristicMove(
  board: Cell[],
  size: BoardSize,
  bot: Player,
  human: Player,
  moves: number[],
): number {
  for (const move of moves) {
    board[move] = bot;
    const won = checkWinner(board, size)?.winner === bot;
    board[move] = null;
    if (won) {
      return move;
    }
  }

  for (const move of moves) {
    board[move] = human;
    const blocked = checkWinner(board, size)?.winner === human;
    board[move] = null;
    if (blocked) {
      return move;
    }
  }

  const center = Math.floor(size / 2) * size + Math.floor(size / 2);
  if (moves.includes(center)) {
    return center;
  }

  let bestScore = Number.NEGATIVE_INFINITY;
  let bestMove = moves[0]!;

  for (const move of moves) {
    board[move] = bot;
    const score = evaluateBoard(board, size, bot, human);
    board[move] = null;

    if (score > bestScore) {
      bestScore = score;
      bestMove = move;
    }
  }

  return bestMove;
}

function evaluateBoard(board: Cell[], size: BoardSize, bot: Player, human: Player): number {
  const winLength = winLengthForSize(size);
  let score = 0;

  const lines = collectLines(size, winLength);
  for (const line of lines) {
    const cells = line.map(index => board[index]!);
    const botCount = cells.filter(cell => cell === bot).length;
    const humanCount = cells.filter(cell => cell === human).length;
    const emptyCount = cells.filter(cell => cell === null).length;

    if (botCount > 0 && humanCount > 0) {
      continue;
    }

    if (botCount > 0) {
      score += 10 ** botCount;
    } else if (humanCount > 0) {
      score -= 10 ** humanCount;
    } else if (emptyCount === winLength) {
      score += 1;
    }
  }

  return score;
}

function collectLines(size: number, winLength: number): number[][] {
  const lines: number[][] = [];

  for (let row = 0; row < size; row++) {
    for (let col = 0; col <= size - winLength; col++) {
      lines.push(Array.from({ length: winLength }, (_, offset) => row * size + col + offset));
    }
  }

  for (let col = 0; col < size; col++) {
    for (let row = 0; row <= size - winLength; row++) {
      lines.push(Array.from({ length: winLength }, (_, offset) => (row + offset) * size + col));
    }
  }

  for (let row = 0; row <= size - winLength; row++) {
    for (let col = 0; col <= size - winLength; col++) {
      lines.push(
        Array.from({ length: winLength }, (_, offset) => (row + offset) * size + (col + offset)),
      );
      lines.push(
        Array.from(
          { length: winLength },
          (_, offset) => (row + offset) * size + (col + winLength - 1 - offset),
        ),
      );
    }
  }

  return lines;
}
