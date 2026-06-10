<script lang="ts" setup>
import type { GameData } from '@/lib/types.ts';
import GameIcon from './GameIcon.vue';

const games :GameData[]= [
    {
        name: 'Minesweeper',
        key: 'minesweeper',
    },
    {
        name: 'Minesweeper',
        key: 'minesweeper',
    },
    {
        name: 'Minesweeper',
        key: 'minesweeper',
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
