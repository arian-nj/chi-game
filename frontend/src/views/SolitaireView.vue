<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useLocalePath } from '@/composables/useLocalePath';
import solitaireLogo from '@/assets/games/solitaire/solitaire_logo.svg';
import SolitaireGameMenu from '@/components/solitaire/SolitaireGameMenu.vue';
import SolitaireSettingsBar from '@/components/solitaire/SolitaireSettingsBar.vue';
import SolitaireHelpSection from '@/components/solitaire/SolitaireHelpSection.vue';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import type { DrawMode } from '@/lib/solitaire/types';
import SolitaireGame from './SolitaireGame.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';

const router = useRouter();
const { localePath } = useLocalePath();

const drawMode = ref<DrawMode>(3);
const gameRef = useTemplateRef('gameRef');

function newGame() {
    gameRef.value?.resetGame();
}

function setDrawMode(mode: DrawMode) {
    drawMode.value = mode;
}

function closeGame() {
    router.push(localePath('/'));
}
</script>

<template>
    <HeaderComponent />

    <div class="solitaire-desktop w-full flex items-start justify-center">
        <XpWindow>
            <XpTitleBar title="Solitaire" :icon="solitaireLogo" @close="closeGame" />
            <SolitaireGameMenu
                :draw-mode="drawMode"
                @new-game="newGame"
                @set-draw-mode="setDrawMode"
            />
            <SolitaireGame ref="gameRef" :draw-mode="drawMode" />
            <SolitaireSettingsBar
                :draw-mode="drawMode"
                @new-game="newGame"
                @set-draw-mode="setDrawMode"
            />
        </XpWindow>
    </div>

    <SolitaireHelpSection />
</template>

<style scoped>
@media (max-width: 640px) {
    .solitaire-desktop {
        align-items: stretch;
        padding: 0.25rem 0.25rem 0;
    }
}
</style>
