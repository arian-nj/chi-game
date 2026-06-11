<script setup lang="ts">
import type { Difficulty } from '@/lib/minesweeper/types';

defineProps<{
    difficulty: Difficulty;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-difficulty': [level: Difficulty];
}>();
</script>

<template>
    <div class="ms-settings">
        <div class="ms-settings-row">
            <span class="ms-settings-label">Difficulty</span>
            <div class="ms-option-group" role="group" aria-label="Difficulty">
                <button
                    type="button"
                    class="ms-option"
                    :class="{ 'ms-option-active': difficulty === 'beginner' }"
                    @click="emit('set-difficulty', 'beginner')"
                >Beginner</button>
                <button
                    type="button"
                    class="ms-option"
                    :class="{ 'ms-option-active': difficulty === 'intermediate' }"
                    @click="emit('set-difficulty', 'intermediate')"
                >Intermediate</button>
                <button
                    type="button"
                    class="ms-option"
                    :class="{ 'ms-option-active': difficulty === 'expert' }"
                    @click="emit('set-difficulty', 'expert')"
                >Expert</button>
            </div>
        </div>

        <div class="ms-settings-actions">
            <button type="button" class="ms-new-game-btn" @click="emit('new-game')">
                New Game
            </button>
        </div>
    </div>
</template>

<style scoped>
.ms-settings {
    background: #c0c0c0;
    padding: 4px 12px 8px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.ms-settings-row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.ms-settings-label {
    color: #000;
    min-width: 52px;
    flex-shrink: 0;
}

.ms-option-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.ms-option {
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

.ms-option:hover:not(.ms-option-active) {
    background: #d4d4d4;
}

.ms-option-active {
    border-color: #808080 #fff #fff #808080;
    background: #ece9d8;
    font-weight: bold;
}

.ms-settings-actions {
    display: flex;
    justify-content: center;
    padding-top: 4px;
    padding-bottom: 2px;
}

.ms-new-game-btn {
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

.ms-new-game-btn:active {
    border-color: #808080 #fff #fff #808080;
}
</style>
