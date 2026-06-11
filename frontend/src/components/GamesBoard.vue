<script lang="ts" setup>
import type { GameData } from '@/lib/types.ts';
import GameIcon from './GameIcon.vue';
import Connect4Logo from '@/assets/games/connect4/connect4_logo.svg';
import MinesweeperLogo from '@/assets/games/minesweeper/minesweeper_logo.png';
import SudokuLogo from '@/assets/games/sudoku/sudoku_logo.svg';
import TicTacToeLogo from '@/assets/games/tictactoe/tictactoe_logo.svg';

const games: GameData[] = [
    {
        name: 'Minesweeper',
        key: 'minesweeper',
        logo: MinesweeperLogo,
    },
    {
        name: 'Tic Tac Toe',
        key: 'tictactoe',
        logo: TicTacToeLogo,
    },
    {
        name: 'Connect 4',
        key: 'connect4',
        logo: Connect4Logo,
    },
    {
        name: 'Sudoku',
        key: 'sudoku',
        logo: SudokuLogo,
    },
]

</script>

<template>
    <div class="grid lg:grid-cols-6 sm:grid-cols-4 grid-cols-3 gap-4
        rounded-2xl">
        <GameIcon v-for="game in games" :key="game.key" :game="game" />
    </div>
</template>

<style scoped>

</style>
