export const gamesData = [
  { key: 'tictactoe' },
  { key: 'conn4' },
  { key: 'linedot' },
  { key: 'go' },
] as const;

export type GameKey = (typeof gamesData)[number]['key'];
