<script setup lang="ts">
import { computed } from 'vue';
import {
    CUSTOM_LIMITS,
    clampCustomConfig,
    maxMinesForBoard,
    type CustomGameConfig,
    type Difficulty,
} from '@/lib/minesweeper/types';

const props = defineProps<{
    difficulty: Difficulty;
    customConfig: CustomGameConfig;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-difficulty': [level: Difficulty];
    'set-custom-config': [config: CustomGameConfig];
}>();

const maxMines = computed(() =>
    maxMinesForBoard(props.customConfig.cellWidth, props.customConfig.cellHeight),
);

function updateCustomField(field: keyof CustomGameConfig, rawValue: string) {
    const parsed = Number.parseInt(rawValue, 10);
    if (Number.isNaN(parsed)) return;

    emit('set-custom-config', clampCustomConfig({
        ...props.customConfig,
        [field]: parsed,
    }));
}
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
                <button
                    type="button"
                    class="ms-option"
                    :class="{ 'ms-option-active': difficulty === 'custom' }"
                    @click="emit('set-difficulty', 'custom')"
                >Custom</button>
            </div>
        </div>

        <div v-if="difficulty === 'custom'" class="ms-custom-panel">
            <div class="ms-custom-field">
                <label class="ms-custom-label" for="ms-custom-width">Width</label>
                <input
                    id="ms-custom-width"
                    class="ms-custom-input"
                    type="number"
                    inputmode="numeric"
                    :min="CUSTOM_LIMITS.minWidth"
                    :max="CUSTOM_LIMITS.maxWidth"
                    :value="customConfig.cellWidth"
                    @change="updateCustomField('cellWidth', ($event.target as HTMLInputElement).value)"
                >
            </div>
            <div class="ms-custom-field">
                <label class="ms-custom-label" for="ms-custom-height">Height</label>
                <input
                    id="ms-custom-height"
                    class="ms-custom-input"
                    type="number"
                    inputmode="numeric"
                    :min="CUSTOM_LIMITS.minHeight"
                    :max="CUSTOM_LIMITS.maxHeight"
                    :value="customConfig.cellHeight"
                    @change="updateCustomField('cellHeight', ($event.target as HTMLInputElement).value)"
                >
            </div>
            <div class="ms-custom-field">
                <label class="ms-custom-label" for="ms-custom-mines">Mines</label>
                <input
                    id="ms-custom-mines"
                    class="ms-custom-input"
                    type="number"
                    inputmode="numeric"
                    :min="CUSTOM_LIMITS.minMines"
                    :max="maxMines"
                    :value="customConfig.mineCount"
                    @change="updateCustomField('mineCount', ($event.target as HTMLInputElement).value)"
                >
            </div>
            <span class="ms-custom-hint">{{ customConfig.cellWidth }}×{{ customConfig.cellHeight }}, max {{ maxMines }} mines</span>
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

.ms-custom-panel {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px 12px;
    padding: 4px 2px 2px;
}

.ms-custom-field {
    display: flex;
    align-items: center;
    gap: 6px;
}

.ms-custom-label {
    color: #000;
    min-width: 38px;
}

.ms-custom-input {
    width: 52px;
    padding: 2px 4px;
    font-size: 11px;
    font-family: inherit;
    border: 2px solid;
    border-color: #808080 #fff #fff #808080;
    background: #fff;
    color: #000;
}

.ms-custom-hint {
    color: #404040;
    font-size: 10px;
    width: 100%;
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
