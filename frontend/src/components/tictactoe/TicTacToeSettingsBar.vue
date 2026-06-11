<script setup lang="ts">
import type { BoardSizeKey, BotDifficulty, GameMode } from '@/lib/tictactoe/types';

defineProps<{
    gameMode: GameMode;
    boardSizeKey: BoardSizeKey;
    botDifficulty: BotDifficulty;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-game-mode': [mode: GameMode];
    'set-board-size': [size: BoardSizeKey];
    'set-bot-difficulty': [level: BotDifficulty];
}>();
</script>

<template>
    <div class="ttt-settings">
        <div class="ttt-settings-row">
            <span class="ttt-settings-label">Mode</span>
            <div class="ttt-option-group" role="group" aria-label="Game mode">
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': gameMode === 'local' }"
                    @click="emit('set-game-mode', 'local')"
                >2 Players</button>
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': gameMode === 'bot' }"
                    @click="emit('set-game-mode', 'bot')"
                >vs Computer</button>
            </div>
        </div>

        <div class="ttt-settings-row">
            <span class="ttt-settings-label">Board</span>
            <div class="ttt-option-group" role="group" aria-label="Board size">
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': boardSizeKey === 'classic' }"
                    @click="emit('set-board-size', 'classic')"
                >3×3</button>
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': boardSizeKey === 'expert' }"
                    @click="emit('set-board-size', 'expert')"
                >5×5</button>
            </div>
        </div>

        <div v-if="gameMode === 'bot'" class="ttt-settings-row">
            <span class="ttt-settings-label">Difficulty</span>
            <div class="ttt-option-group" role="group" aria-label="Bot difficulty">
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': botDifficulty === 'easy' }"
                    @click="emit('set-bot-difficulty', 'easy')"
                >Easy</button>
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': botDifficulty === 'medium' }"
                    @click="emit('set-bot-difficulty', 'medium')"
                >Medium</button>
                <button
                    type="button"
                    class="ttt-option"
                    :class="{ 'ttt-option-active': botDifficulty === 'hard' }"
                    @click="emit('set-bot-difficulty', 'hard')"
                >Hard</button>
            </div>
        </div>

        <div class="ttt-settings-actions">
            <button type="button" class="ttt-new-game-btn" @click="emit('new-game')">
                New Game
            </button>
        </div>
    </div>
</template>

<style scoped>
.ttt-settings {
    background: #c0c0c0;
    padding: 4px 12px 8px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.ttt-settings-row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.ttt-settings-label {
    color: #000;
    min-width: 52px;
    flex-shrink: 0;
}

.ttt-option-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.ttt-option {
    padding: 3px 10px;
    font-size: 11px;
    font-family: inherit;
    line-height: 1.2;
    color: #000;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
    cursor: default;
    white-space: nowrap;
}

.ttt-option:hover:not(.ttt-option-active) {
    background: #d4d4d4;
}

.ttt-option-active {
    border-color: #808080 #fff #fff #808080;
    background: #ece9d8;
    font-weight: bold;
}

.ttt-settings-actions {
    display: flex;
    justify-content: center;
    padding-top: 4px;
    padding-bottom: 2px;
}

.ttt-new-game-btn {
    padding: 4px 16px;
    font-size: 11px;
    font-family: inherit;
    line-height: 1.2;
    color: #000;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
    cursor: default;
}

.ttt-new-game-btn:active {
    border-color: #808080 #fff #fff #808080;
}
</style>
