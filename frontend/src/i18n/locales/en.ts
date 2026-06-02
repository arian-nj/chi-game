export default {
  app: {
    title: 'Chi Game',
  },
  seo: {
    siteName: 'Chi Game',
    defaultDescription:
      'Play free board games in your browser — Tic Tac Toe, Connect 4, and more. Offline-friendly, vs bot or two players.',
    home: {
      title: 'Chi Game — Play Tic Tac Toe & Connect 4 Online',
      description:
        'Free browser games at Chi Game. Play Tic Tac Toe and Connect 4 against a bot or a friend — no install required.',
    },
    changelog: {
      title: 'Changelog | Chi Game',
      description: 'What is new, changed, and fixed in Chi Game.',
    },
    about: {
      title: 'About | Chi Game',
      description:
        'Learn about Chi Game — free browser board games, offline play, bot and two-player modes, and what is coming next.',
    },
    notFound: {
      title: 'Page not found | Chi Game',
      description: "This page doesn't exist. Head back to Chi Game and pick a game to play.",
    },
    gamePlay: {
      tictactoe: {
        title: 'Play Tic Tac Toe | Chi Game',
        description:
          'Play Tic Tac Toe (XO) free in your browser. Choose 3×3 or 5×5, play vs bot or two players.',
      },
      conn4: {
        title: 'Play Connect 4 | Chi Game',
        description:
          'Play Connect 4 free online. Drop discs on a classic 7×6 board — vs bot or two players.',
      },
      fallback: {
        title: 'Play {game} | Chi Game',
        description: 'Play {game} free in your browser at Chi Game.',
      },
    },
    gameRules: {
      tictactoe: {
        title: 'Tic Tac Toe Rules | Chi Game',
        description:
          'Learn how to play Tic Tac Toe: objective, turns, winning lines, draws, and a quick strategy tip.',
      },
      conn4: {
        title: 'Connect 4 Rules | Chi Game',
        description:
          'Learn how to play Connect 4: board setup, turns, four-in-a-row wins, draws, and strategy.',
      },
      fallback: {
        title: '{game} Rules | Chi Game',
        description: 'How to play {game} — rules and tips at Chi Game.',
      },
    },
  },
  notFound: {
    title: 'Page not found',
    description: "This page doesn't exist (or was moved).",
    backHome: 'Back to home',
  },
  nav: {
    home: 'Home',
    play: 'Play',
    rules: 'Rules',
    backToHome: 'Back to home',
    changelog: 'Changelog',
    about: 'About',
  },
  about: {
    title: 'About Chi Game',
    tagline:
      'Chi Game is a free collection of classic board games you can play in your browser — no install required.',
    sections: {
      mission: {
        title: 'What we are building',
        body: 'We want a simple place to pick a game, learn the rules, and start playing right away. Everything runs in the browser and works offline once loaded.',
      },
      games: {
        title: 'Games today',
        body: 'Phase 1 focuses on single-device play: challenge the bot or pass the device with a friend. Tic Tac Toe and Connect 4 are available now, with more classics on the way.',
      },
      roadmap: {
        title: 'What is next',
        body: 'Later phases will add online multiplayer so you can play with friends remotely. Until then, enjoy quick offline matches and check the changelog for updates.',
      },
    },
  },
  changelog: {
    title: 'Changelog',
    sections: {
      new: 'New',
      changed: 'Changed',
      fixed: 'Fixed',
    },
  },
  locale: {
    label: 'Language',
    en: 'English',
    fa: 'Persian',
  },
  games: {
    tictactoe: 'Tic Tac Toe',
    conn4: 'Connect 4',
    linedot: 'Line Dot',
    go: 'Go!',
  },
  settings: {
    opponent: 'Opponent',
    bot: 'Bot 🤖',
    twoPlayer: '👥 2 Player',
    boardSize: 'Board size',
    board3x3: '3×3',
    board5x5: '5×5',
    connect4Hint: 'Classic 7×6 board — connect four discs in a row to win.',
    difficulty: 'Bot difficulty',
    difficultyEasy: 'Easy',
    difficultyMedium: 'Medium',
    difficultyHard: 'Hard',
  },
  play: {
    button: 'Play',
  },
  game: {
    newGame: 'New game',
    backToSettings: 'Back to settings',
    youWin: 'You win!',
    botWins: 'Bot wins!',
    draw: "It's a draw!",
    botThinking: 'Bot is thinking…',
    playerWins: 'Player {player} wins!',
    yourTurnX: 'Your turn (X)',
    botTurnO: 'Bot turn (O)',
    playerTurn: "Player {player}'s turn",
    redWins: 'Red wins!',
    yellowWins: 'Yellow wins!',
    yourTurnRed: 'Your turn (Red)',
    botTurnYellow: 'Bot turn (Yellow)',
    redTurn: 'Red turn',
    yellowTurn: 'Yellow turn',
    red: 'Red',
    yellow: 'Yellow',
    you: 'You',
    botPlayer: 'Bot',
    playerX: 'Player X',
    playerO: 'Player O',
    tttBoardAria: '{size} by {size} tic tac toe board',
    cellOccupied: 'Cell {index}, {symbol}',
    cellEmpty: 'Empty cell {index}',
    connect4BoardAria: 'Connect 4 board',
    dropDiscColumn: 'Drop disc in column {col}',
    columnRow: 'Column {col}, row {row}',
  },
  rules: {
    tictactoe: {
      title: 'How to Play XO (Tic Tac Toe)',
      objective: 'Objective',
      objectiveText:
        'The goal of the game is to be the first player to align {highlight} in a continuous straight line. This line can be horizontal, vertical, or diagonal.',
      objectiveHighlight: 'three of your symbols',
      gameplay: 'Gameplay Rules',
      rule1: 'The match is played on a classic {highlight} consisting of 9 empty squares.',
      rule1Highlight: '3x3 grid',
      rule2: 'Player 1 is assigned the {x} symbol, and Player 2 (or the AI) takes the {o} symbol.',
      rule3: 'Players take turns placing their respective marks into any remaining unoccupied square.',
      rule4: 'Once a mark is placed, it cannot be moved, overwritten, or removed from the board.',
      winning: 'Winning & Draws',
      victory: '{label} The instant a player successfully arranges 3 marks in a row, column, or diagonal, the game ends and that player wins.',
      victoryLabel: 'Victory:',
      drawRule:
        '{label} If all 9 squares on the grid are completely filled and neither player has achieved a line of 3, the game ends in a stalemate/tie.',
      drawLabel: 'Draw (Cats Game):',
      hintTitle: 'Quick Strategy Hint:',
      hintText:
        'If you go first, controlling the center square or the outer corners gives you the highest mathematical probability to force a winning trap!',
    },
    conn4: {
      title: 'How to Play Connect 4',
      objective: 'Objective',
      objectiveText:
        'Be the first player to connect {highlight} in a straight line. The line can be horizontal, vertical, or diagonal.',
      objectiveHighlight: 'four of your discs',
      gameplay: 'Gameplay Rules',
      rule1: 'The game is played on a {highlight} with 7 columns and 6 rows.',
      rule1Highlight: '7×6 vertical board',
      rule2: 'Player 1 uses {red} discs, and Player 2 (or the bot) uses {yellow} discs.',
      rule3: 'On your turn, pick a column; your disc drops into the lowest empty slot in that column.',
      rule4: 'Once a disc is placed, it cannot be moved. You cannot play in a full column.',
      winning: 'Winning & Draws',
      victory:
        '{label} Connect four of your discs in a row, column, or diagonal to win immediately.',
      victoryLabel: 'Victory:',
      drawRule:
        '{label} If every slot on the board is filled and neither player has four in a row, the game ends in a draw.',
      drawLabel: 'Draw:',
      hintTitle: 'Quick Strategy Hint:',
      hintText:
        'Watch the center columns, set up double threats, and block your opponent when they are one move from four in a row.',
    },
  },
} as const;
