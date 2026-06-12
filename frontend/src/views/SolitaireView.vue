<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useLocalePath } from '@/composables/useLocalePath';
import solitaireLogo from '@/assets/games/solitaire/solitaire_logo.svg';
import XpTitleBar from '@/components/xp/XpTitleBar.vue';
import XpWindow from '@/components/xp/XpWindow.vue';
import HeaderComponent from '@/components/header/HeaderComponent.vue';
import SolitaireGame from './SolitaireGame.vue';

const router = useRouter();
const { localePath } = useLocalePath();

const drawMode = ref<1 | 3>(1);
const gameRef = useTemplateRef('gameRef');

function newGame() {
    gameRef.value?.resetGame();
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
            <div class="sol-menu-bar">
                <button type="button" class="sol-menu-btn" @click="newGame">New Game</button>
            </div>
            <SolitaireGame ref="gameRef" :draw-mode="drawMode" />
        </XpWindow>
    </div>
</template>

<style scoped>
.sol-menu-bar {
    display: flex;
    gap: 4px;
    padding: 4px 6px;
    background: #ece9d8;
    border-bottom: 1px solid #aca899;
}

.sol-menu-btn {
    padding: 2px 10px;
    font-size: 12px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    background: #ece9d8;
    border: 1px solid transparent;
    cursor: default;
}

.sol-menu-btn:hover {
    border-color: #aca899;
    background: #fff;
}

@media (max-width: 640px) {
    .solitaire-desktop {
        align-items: stretch;
        padding: 0.25rem 0.25rem 0;
    }
}
</style>
