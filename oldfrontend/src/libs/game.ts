export const gamesData = [
  { key: 'tic-tac-toe',name: 'Tic Tac Toe',isEnable:true },
  { key: 'connect-4',name: 'Connect 4',isEnable:true },
  { key: 'line-dot',name: 'Line Dot',isEnable:false },
  { key: 'go',name: 'Go',isEnable:false },
] as const;

export type GameKey = (typeof gamesData)[number]['key'];
