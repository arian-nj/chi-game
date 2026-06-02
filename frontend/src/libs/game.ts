export const gamesData = [
  { key: 'tictactoe',isEnable:true },
  { key: 'conn4',isEnable:true },
  { key: 'linedot',isEnable:false },
  { key: 'go',isEnable:false },
] as const;

export type GameKey = (typeof gamesData)[number]['key'];
