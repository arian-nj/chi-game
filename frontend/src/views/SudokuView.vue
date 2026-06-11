<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import sudokuLogo from '@/assets/games/sudoku/sudoku_logo.svg';
import SudokuGameMenu from '@/components/sudoku/SudokuGameMenu.vue';
import SudokuSettingsBar from '@/components/sudoku/SudokuSettingsBar.vue';
import SudokuHelpSection from '@/components/sudoku/SudokuHelpSection.vue';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import type { Difficulty } from '@/lib/sudoku/types';
import SudokuGame from './SudokuGame.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';

const router = useRouter();

const difficulty = ref<Difficulty>('easy');
const gameRef = useTemplateRef('gameRef');

function newGame() {
    gameRef.value?.resetGame();
}

function setDifficulty(level: Difficulty) {
    difficulty.value = level;
}

function closeGame() {
    router.push('/');
}
</script>

<template>
    <HeaderComponent />

    <div class="sudoku-desktop min-h-screen w-full flex items-start justify-center">
        <XpWindow>
            <XpTitleBar title="Sudoku" :icon="sudokuLogo" @close="closeGame" />
            <SudokuGameMenu
                :difficulty="difficulty"
                @new-game="newGame"
                @set-difficulty="setDifficulty"
            />
            <SudokuGame
                ref="gameRef"
                :difficulty="difficulty"
            />
            <SudokuSettingsBar
                :difficulty="difficulty"
                @new-game="newGame"
                @set-difficulty="setDifficulty"
            />
        </XpWindow>
    </div>

    <SudokuHelpSection />
</template>

<style scoped>
@media (max-width: 640px) {
    .sudoku-desktop {
        padding: 0.5rem 0.25rem;
    }
}
</style>
