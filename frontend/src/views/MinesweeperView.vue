<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, useTemplateRef } from 'vue';
import MinesweeperGame from './MinesweeperGame.vue';
import { useRouter } from 'vue-router';
import minesweeperLogo from '@/assets/games/minesweeper/minesweeper_logo.png';

const router = useRouter();

type Difficulty = 'beginner' | 'intermediate' | 'expert';

const DIFFICULTIES: Record<Difficulty, { cellWidth: number; cellHeight: number; mineCount: number }> = {
    beginner: { cellWidth: 9, cellHeight: 9, mineCount: 10 },
    intermediate: { cellWidth: 16, cellHeight: 16, mineCount: 40 },
    expert: { cellWidth: 30, cellHeight: 16, mineCount: 99 },
};

const difficulty = ref<Difficulty>('beginner');
const gameMenuOpen = ref(false);
const gameRef = useTemplateRef('gameRef');
const menuRootRef = useTemplateRef('menuRootRef');

const gameConfig = computed(() => DIFFICULTIES[difficulty.value]);

function closeGameMenu() {
    gameMenuOpen.value = false;
}

function toggleGameMenu() {
    gameMenuOpen.value = !gameMenuOpen.value;
}

function newGame() {
    gameRef.value?.resetGame();
    closeGameMenu();
}

function setDifficulty(level: Difficulty) {
    difficulty.value = level;
    closeGameMenu();
}

function onDocumentClick(event: MouseEvent) {
    if (!gameMenuOpen.value) return;
    const root = menuRootRef.value;
    if (root && !root.contains(event.target as Node)) {
        closeGameMenu();
    }
}
function closeGame() {
    router.push('/');
}

onMounted(() => document.addEventListener('click', onDocumentClick));
onUnmounted(() => document.removeEventListener('click', onDocumentClick));
</script>

<template>
    <div class="minesweeper-desktop min-h-screen w-full flex items-start justify-center">
        <div class="xp-window">
            <!-- title bar -->
            <div class="xp-titlebar">
                <div class="xp-titlebar-left">
                    <img
                        :src="minesweeperLogo"
                        alt="Minesweeper Logo"
                        class="xp-titlebar-icon"
                    >
                    <span class="xp-titlebar-text">Minesweeper</span>
                </div>
                <div class="xp-titlebar-buttons">
                    <button type="button" class="xp-win-btn xp-win-btn-min" aria-label="Minimize" />
                    <button type="button" class="xp-win-btn xp-win-btn-max" aria-label="Maximize" />
                    <button type="button" @click="closeGame" class="xp-win-btn xp-win-btn-close" aria-label="Close" />
                </div>
            </div>

            <!-- menu bar -->
            <div class="xp-menubar">
                <div ref="menuRootRef" class="xp-menu-root">
                    <button
                        type="button"
                        class="xp-menu-item"
                        :class="{ 'xp-menu-item-active': gameMenuOpen }"
                        @click.stop="toggleGameMenu"
                    >Game</button>
                    <div v-if="gameMenuOpen" class="xp-dropdown" @click.stop>
                        <button type="button" class="xp-dropdown-item" @click="newGame">
                            <span class="xp-dropdown-bullet" />
                            New Game
                        </button>
                        <div class="xp-dropdown-separator" />
                        <button
                            type="button"
                            class="xp-dropdown-item"
                            @click="setDifficulty('beginner')"
                        >
                            <span class="xp-dropdown-bullet">{{ difficulty === 'beginner' ? '•' : '' }}</span>
                            Beginner
                        </button>
                        <button
                            type="button"
                            class="xp-dropdown-item"
                            @click="setDifficulty('intermediate')"
                        >
                            <span class="xp-dropdown-bullet">{{ difficulty === 'intermediate' ? '•' : '' }}</span>
                            Intermediate
                        </button>
                        <button
                            type="button"
                            class="xp-dropdown-item"
                            @click="setDifficulty('expert')"
                        >
                            <span class="xp-dropdown-bullet">{{ difficulty === 'expert' ? '•' : '' }}</span>
                            Expert
                        </button>
                    </div>
                </div>
                <a href="#how-to-play-minesweeper" class="xp-menu-item">Help</a>
            </div>

            <!-- game -->
            <MinesweeperGame
                ref="gameRef"
                :cell-height="gameConfig.cellHeight"
                :cell-width="gameConfig.cellWidth"
                :mine-count="gameConfig.mineCount"
            />
        </div>
    </div>

    <div class="bg-gray-400" id="how-to-play-minesweeper">
        <h1>How to play minesweeper</h1>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
    </div>
</template>

<style scoped>
.xp-window {
    --title-btn-size: 28px;
    --title-icon-size: 18px;
    --title-font-size: 13px;
    border: 1px solid #0054e3;
    border-radius: 8px 8px 0 0;
    box-shadow: 1px 1px 0 #000;
    overflow: hidden;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    width: fit-content;
    max-width: 100%;
}

@media (max-width: 640px) {
    .minesweeper-desktop {
        padding: 0.5rem 0.25rem;
    }

    .xp-window {
        --title-btn-size: 34px;
        --title-icon-size: 22px;
        --title-font-size: 16px;
    }
}

.xp-titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    min-height: calc(var(--title-btn-size) + 8px);
    padding: 4px 4px 4px 5px;
    background: linear-gradient(180deg, #0997ff 0%, #0053ee 8%, #0050ee 40%, #06f 88%, #06f 93%, #005eff 95%, #003ddb 96%, #003ddb 100%);
    border-radius: 8px 8px 0 0;
}

.xp-titlebar-left {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
}

.xp-titlebar-icon {
    width: var(--title-icon-size);
    height: var(--title-icon-size);
    flex-shrink: 0;
}

.xp-titlebar-text {
    color: #fff;
    font-size: var(--title-font-size);
    font-weight: bold;
    text-shadow: 1px 1px #000;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1;
}

.xp-titlebar-buttons {
    display: flex;
    align-items: center;
    gap: 2px;
    flex-shrink: 0;
}

.xp-win-btn {
    appearance: none;
    -webkit-appearance: none;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: var(--title-btn-size);
    height: var(--title-btn-size);
    min-width: var(--title-btn-size);
    max-width: var(--title-btn-size);
    box-sizing: border-box;
    flex-shrink: 0;
    padding: 0;
    margin: 0;
    cursor: default;
    border-radius: 3px;
}

.xp-win-btn-min,
.xp-win-btn-max {
    background: linear-gradient(180deg, #3c8cfd 0%, #3b8bfc 50%, #3a8afb 51%, #3586f8 100%);
    border: 1px solid #003c74;
    box-shadow: inset 1px 1px 0 #6ba8ff;
}

.xp-win-btn-min::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 55%;
    height: 2px;
    background: #fff;
}

.xp-win-btn-max::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 42%;
    height: 42%;
    border: 2px solid #fff;
    box-sizing: border-box;
}

.xp-win-btn-close {
    background: linear-gradient(180deg, #e88a8f 0%, #e4787d 50%, #e36b70 51%, #d94a4f 100%);
    border: 1px solid #8b0000;
    box-shadow: inset 1px 1px 0 #f0a8ab;
}

.xp-win-btn-close::before,
.xp-win-btn-close::after {
    content: '';
    position: absolute;
    left: 50%;
    top: 50%;
    width: 55%;
    height: 2px;
    background: #fff;
}

.xp-win-btn-close::before {
    transform: translate(-50%, -50%) rotate(45deg);
}

.xp-win-btn-close::after {
    transform: translate(-50%, -50%) rotate(-45deg);
}

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
