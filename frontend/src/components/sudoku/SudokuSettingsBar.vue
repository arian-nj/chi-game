<script setup lang="ts">
import type { Difficulty } from '@/lib/sudoku/types';

defineProps<{
    difficulty: Difficulty;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-difficulty': [level: Difficulty];
}>();
</script>

<template>
    <div class="sd-settings">
        <div class="sd-settings-row">
            <span class="sd-settings-label">Difficulty</span>
            <div class="sd-option-group" role="group" aria-label="Difficulty">
                <button
                    type="button"
                    class="sd-option"
                    :class="{ 'sd-option-active': difficulty === 'easy' }"
                    @click="emit('set-difficulty', 'easy')"
                >Easy</button>
                <button
                    type="button"
                    class="sd-option"
                    :class="{ 'sd-option-active': difficulty === 'medium' }"
                    @click="emit('set-difficulty', 'medium')"
                >Medium</button>
                <button
                    type="button"
                    class="sd-option"
                    :class="{ 'sd-option-active': difficulty === 'hard' }"
                    @click="emit('set-difficulty', 'hard')"
                >Hard</button>
            </div>
        </div>

        <div class="sd-settings-actions">
            <button type="button" class="sd-new-game-btn" @click="emit('new-game')">
                New Game
            </button>
        </div>
    </div>
</template>

<style scoped>
.sd-settings {
    background: #c0c0c0;
    padding: 4px 12px 8px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.sd-settings-row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.sd-settings-label {
    color: #000;
    min-width: 52px;
    flex-shrink: 0;
}

.sd-option-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.sd-option {
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

.sd-option:hover:not(.sd-option-active) {
    background: #d4d4d4;
}

.sd-option-active {
    border-color: #808080 #fff #fff #808080;
    background: #ece9d8;
    font-weight: bold;
}

.sd-settings-actions {
    display: flex;
    justify-content: center;
    padding-top: 4px;
    padding-bottom: 2px;
}

.sd-new-game-btn {
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

.sd-new-game-btn:active {
    border-color: #808080 #fff #fff #808080;
}
</style>
