<script lang="ts" setup>
import type { GameData } from '@/lib/types.ts';
import GameIcon from './GameIcon.vue';
import MinesweeperLogo from '@/assets/games/minesweeper/minesweeper_logo.png';
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
]

</script>

<template>
    <div class="grid lg:grid-cols-6 sm:grid-cols-4 grid-cols-3 gap-4 mt-8
        rounded-2xl">
        <GameIcon v-for="game in games" :key="game.key" :game="game" />
    </div>
</template>

<style scoped>

</style>
