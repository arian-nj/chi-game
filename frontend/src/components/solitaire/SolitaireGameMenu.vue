<script setup lang="ts">
import { onMounted, onUnmounted, ref, useTemplateRef } from 'vue';
import type { DrawMode } from '@/lib/solitaire/types';

defineProps<{
    drawMode: DrawMode;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-draw-mode': [mode: DrawMode];
}>();

const gameMenuOpen = ref(false);
const menuRootRef = useTemplateRef('menuRootRef');

function closeGameMenu() {
    gameMenuOpen.value = false;
}

function toggleGameMenu() {
    gameMenuOpen.value = !gameMenuOpen.value;
}

function onNewGame() {
    emit('new-game');
    closeGameMenu();
}

function onSetDrawMode(mode: DrawMode) {
    emit('set-draw-mode', mode);
    closeGameMenu();
}

function onDocumentClick(event: MouseEvent) {
    if (!gameMenuOpen.value) return;
    const root = menuRootRef.value;
    if (root && !root.contains(event.target as Node)) {
        closeGameMenu();
    }
}

onMounted(() => document.addEventListener('click', onDocumentClick));
onUnmounted(() => document.removeEventListener('click', onDocumentClick));
</script>

<template>
    <div class="xp-menubar">
        <div ref="menuRootRef" class="xp-menu-root">
            <button
                type="button"
                class="xp-menu-item"
                :class="{ 'xp-menu-item-active': gameMenuOpen }"
                @click.stop="toggleGameMenu"
            >Game</button>
            <div v-if="gameMenuOpen" class="xp-dropdown" @click.stop>
                <button type="button" class="xp-dropdown-item" @click="onNewGame">
                    <span class="xp-dropdown-bullet" />
                    New Game
                </button>
                <div class="xp-dropdown-separator" />
                <button type="button" class="xp-dropdown-item" @click="onSetDrawMode(1)">
                    <span class="xp-dropdown-bullet">{{ drawMode === 1 ? '•' : '' }}</span>
                    Draw 1
                </button>
                <button type="button" class="xp-dropdown-item" @click="onSetDrawMode(3)">
                    <span class="xp-dropdown-bullet">{{ drawMode === 3 ? '•' : '' }}</span>
                    Draw 3
                </button>
            </div>
        </div>
        <a href="#how-to-play-solitaire" class="xp-menu-item">Help</a>
    </div>
</template>

<style scoped>
.xp-menubar {
    display: flex;
    align-items: center;
    gap: 0;
    background: #ece9d8;
    padding: 2px 0;
    border-bottom: 1px solid #aca899;
    position: relative;
    z-index: 10;
}

.xp-menu-root {
    position: relative;
    display: flex;
    align-items: center;
}

.xp-menu-item {
    display: inline-flex;
    align-items: center;
    margin: 0;
    background: none;
    border: none;
    color: #000;
    font-size: 15px;
    line-height: 1;
    font-family: inherit;
    padding: 2px 6px;
    cursor: default;
    text-decoration: none;
    box-sizing: border-box;
}

.xp-menu-item:hover,
.xp-menu-item-active {
    background: #316ac5;
    color: #fff;
}

.xp-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    min-width: 140px;
    background: #fff;
    border: 1px solid #716f64;
    box-shadow: 2px 2px 2px rgba(0, 0, 0, 0.2);
    padding: 2px 0;
    z-index: 20;
}

.xp-dropdown-item {
    display: flex;
    align-items: center;
    width: 100%;
    background: none;
    border: none;
    color: #000;
    font-size: 11px;
    font-family: inherit;
    padding: 3px 20px 3px 2px;
    cursor: default;
    text-align: left;
    white-space: nowrap;
}

.xp-dropdown-item:hover {
    background: #316ac5;
    color: #fff;
}

.xp-dropdown-bullet {
    display: inline-block;
    width: 14px;
    text-align: center;
    flex-shrink: 0;
}

.xp-dropdown-separator {
    height: 1px;
    margin: 2px 2px;
    background: #aca899;
}
</style>
