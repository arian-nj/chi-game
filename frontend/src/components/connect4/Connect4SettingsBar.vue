<script setup lang="ts">
import type { BotDifficulty, GameMode } from '@/lib/connect4/types';

defineProps<{
    gameMode: GameMode;
    botDifficulty: BotDifficulty;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-game-mode': [mode: GameMode];
    'set-bot-difficulty': [level: BotDifficulty];
}>();
</script>

<template>
    <div class="c4-settings">
        <div class="c4-settings-row">
            <span class="c4-settings-label">Mode</span>
            <div class="c4-option-group" role="group" aria-label="Game mode">
                <button
                    type="button"
                    class="c4-option"
                    :class="{ 'c4-option-active': gameMode === 'local' }"
                    @click="emit('set-game-mode', 'local')"
                >2 Players</button>
                <button
                    type="button"
                    class="c4-option"
                    :class="{ 'c4-option-active': gameMode === 'bot' }"
                    @click="emit('set-game-mode', 'bot')"
                >vs Computer</button>
            </div>
        </div>

        <div v-if="gameMode === 'bot'" class="c4-settings-row">
            <span class="c4-settings-label">Difficulty</span>
            <div class="c4-option-group" role="group" aria-label="Bot difficulty">
                <button
                    type="button"
                    class="c4-option"
                    :class="{ 'c4-option-active': botDifficulty === 'easy' }"
                    @click="emit('set-bot-difficulty', 'easy')"
                >Easy</button>
                <button
                    type="button"
                    class="c4-option"
                    :class="{ 'c4-option-active': botDifficulty === 'medium' }"
                    @click="emit('set-bot-difficulty', 'medium')"
                >Medium</button>
                <button
                    type="button"
                    class="c4-option"
                    :class="{ 'c4-option-active': botDifficulty === 'hard' }"
                    @click="emit('set-bot-difficulty', 'hard')"
                >Hard</button>
            </div>
        </div>

        <div class="c4-settings-actions">
            <button type="button" class="c4-new-game-btn" @click="emit('new-game')">
                New Game
            </button>
        </div>
    </div>
</template>

<style scoped>
.c4-settings {
    background: #c0c0c0;
    padding: 4px 12px 8px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.c4-settings-row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.c4-settings-label {
    color: #000;
    min-width: 52px;
    flex-shrink: 0;
}

.c4-option-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.c4-option {
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

.c4-option:hover:not(.c4-option-active) {
    background: #d4d4d4;
}

.c4-option-active {
    border-color: #808080 #fff #fff #808080;
    background: #ece9d8;
    font-weight: bold;
}

.c4-settings-actions {
    display: flex;
    justify-content: center;
    padding-top: 4px;
    padding-bottom: 2px;
}

.c4-new-game-btn {
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

.c4-new-game-btn:active {
    border-color: #808080 #fff #fff #808080;
}
</style>
