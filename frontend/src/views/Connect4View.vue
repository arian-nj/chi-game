<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useLocalePath } from '@/composables/useLocalePath';
import connect4Logo from '@/assets/games/connect4/connect4_logo.svg';
import Connect4GameMenu from '@/components/connect4/Connect4GameMenu.vue';
import Connect4SettingsBar from '@/components/connect4/Connect4SettingsBar.vue';
import Connect4HelpSection from '@/components/connect4/Connect4HelpSection.vue';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import type { BotDifficulty, GameMode } from '@/lib/connect4/types';
import Connect4Game from './Connect4Game.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';

const router = useRouter();
const { localePath } = useLocalePath();

const gameMode = ref<GameMode>('bot');
const botDifficulty = ref<BotDifficulty>('medium');
const gameRef = useTemplateRef('gameRef');

const isBot = computed(() => gameMode.value === 'bot');

function newGame() {
    gameRef.value?.resetGame();
}

function setGameMode(mode: GameMode) {
    gameMode.value = mode;
}

function setBotDifficulty(level: BotDifficulty) {
    botDifficulty.value = level;
}

function closeGame() {
    router.push(localePath('/'));
}
</script>

<template>
    <HeaderComponent />

    <div class="connect4-desktop w-full flex items-start justify-center">
        <XpWindow>
            <XpTitleBar title="Connect 4" :icon="connect4Logo" @close="closeGame" />
            <Connect4GameMenu
                :game-mode="gameMode"
                :bot-difficulty="botDifficulty"
                @new-game="newGame"
                @set-game-mode="setGameMode"
                @set-bot-difficulty="setBotDifficulty"
            />
            <Connect4Game
                ref="gameRef"
                :is-bot="isBot"
                :bot-difficulty="botDifficulty"
            />
            <Connect4SettingsBar
                :game-mode="gameMode"
                :bot-difficulty="botDifficulty"
                @new-game="newGame"
                @set-game-mode="setGameMode"
                @set-bot-difficulty="setBotDifficulty"
            />
        </XpWindow>
    </div>

    <Connect4HelpSection />
</template>

<style scoped>
@media (max-width: 640px) {
    .connect4-desktop {
        padding: 0.5rem 0.25rem;
    }
}
</style>
