<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useLocalePath } from '@/composables/useLocalePath';
import minesweeperLogo from '@/assets/games/minesweeper/minesweeper_logo.png';
import MinesweeperGameMenu from '@/components/minesweeper/MinesweeperGameMenu.vue';
import MinesweeperSettingsBar from '@/components/minesweeper/MinesweeperSettingsBar.vue';
import MinesweeperHelpSection from '@/components/minesweeper/MinesweeperHelpSection.vue';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import { DIFFICULTIES, DEFAULT_CUSTOM_CONFIG, clampCustomConfig, type CustomGameConfig, type Difficulty } from '@/lib/minesweeper/types';
import MinesweeperGame from './MinesweeperGame.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';

const router = useRouter();
const { localePath } = useLocalePath();

const difficulty = ref<Difficulty>('beginner');
const customConfig = ref<CustomGameConfig>({ ...DEFAULT_CUSTOM_CONFIG });
const gameRef = useTemplateRef('gameRef');

const gameConfig = computed(() => {
    if (difficulty.value === 'custom') {
        return clampCustomConfig(customConfig.value);
    }
    return DIFFICULTIES[difficulty.value];
});

function newGame() {
    gameRef.value?.resetGame();
}

function setDifficulty(level: Difficulty) {
    difficulty.value = level;
}

function setCustomConfig(config: CustomGameConfig) {
    customConfig.value = clampCustomConfig(config);
    difficulty.value = 'custom';
}

function closeGame() {
    router.push(localePath('/'));
}
</script>

<template>
    <HeaderComponent />
    <div class="minesweeper-desktop w-full flex items-start justify-center">
        <XpWindow>
            <XpTitleBar title="Minesweeper" :icon="minesweeperLogo" @close="closeGame" />
            <MinesweeperGameMenu
                :difficulty="difficulty"
                @new-game="newGame"
                @set-difficulty="setDifficulty"
            />
            <MinesweeperGame
                ref="gameRef"
                :cell-height="gameConfig.cellHeight"
                :cell-width="gameConfig.cellWidth"
                :mine-count="gameConfig.mineCount"
            />
            <MinesweeperSettingsBar
                :difficulty="difficulty"
                :custom-config="customConfig"
                @new-game="newGame"
                @set-difficulty="setDifficulty"
                @set-custom-config="setCustomConfig"
            />
        </XpWindow>
    </div>

    <MinesweeperHelpSection />
</template>

<style scoped>
@media (max-width: 640px) {
    .minesweeper-desktop {
        padding: 0.5rem 0.25rem;
    }
}
</style>
