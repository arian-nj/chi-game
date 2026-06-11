<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import tictactoeLogo from '@/assets/games/tictactoe/tictactoe_logo.svg';
import TicTacToeGameMenu from '@/components/tictactoe/TicTacToeGameMenu.vue';
import TicTacToeSettingsBar from '@/components/tictactoe/TicTacToeSettingsBar.vue';
import TicTacToeHelpSection from '@/components/tictactoe/TicTacToeHelpSection.vue';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import { BOARD_SIZES, type BoardSizeKey, type BotDifficulty, type GameMode } from '@/lib/tictactoe/types';
import TicTacToeGame from './TicTacToeGame.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';

const router = useRouter();

const gameMode = ref<GameMode>('bot');
const boardSizeKey = ref<BoardSizeKey>('classic');
const botDifficulty = ref<BotDifficulty>('medium');
const gameRef = useTemplateRef('gameRef');

const boardSize = computed(() => BOARD_SIZES[boardSizeKey.value]);
const isBot = computed(() => gameMode.value === 'bot');

function newGame() {
    gameRef.value?.resetGame();
}

function setGameMode(mode: GameMode) {
    gameMode.value = mode;
}

function setBoardSize(size: BoardSizeKey) {
    boardSizeKey.value = size;
}

function setBotDifficulty(level: BotDifficulty) {
    botDifficulty.value = level;
}

function closeGame() {
    router.push('/');
}
</script>

<template>
    <HeaderComponent />

    <div class="tictactoe-desktop w-full flex items-start justify-center">
        <XpWindow>
            <XpTitleBar title="Tic Tac Toe" :icon="tictactoeLogo" @close="closeGame" />
            <TicTacToeGameMenu
                :game-mode="gameMode"
                :board-size-key="boardSizeKey"
                :bot-difficulty="botDifficulty"
                @new-game="newGame"
                @set-game-mode="setGameMode"
                @set-board-size="setBoardSize"
                @set-bot-difficulty="setBotDifficulty"
            />
            <TicTacToeGame
                ref="gameRef"
                :board-size="boardSize"
                :is-bot="isBot"
                :bot-difficulty="botDifficulty"
            />
            <TicTacToeSettingsBar
                :game-mode="gameMode"
                :board-size-key="boardSizeKey"
                :bot-difficulty="botDifficulty"
                @new-game="newGame"
                @set-game-mode="setGameMode"
                @set-board-size="setBoardSize"
                @set-bot-difficulty="setBotDifficulty"
            />
        </XpWindow>
    </div>

    <TicTacToeHelpSection />
</template>

<style scoped>
@media (max-width: 640px) {
    .tictactoe-desktop {
        padding: 0.5rem 0.25rem;
    }
}
</style>
