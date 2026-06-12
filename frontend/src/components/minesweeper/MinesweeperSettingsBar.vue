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

function updateCustomField(field: keyof CustomGameConfig, rawValue: number) {
    if (Number.isNaN(rawValue)) return;

    emit('set-custom-config', clampCustomConfig({
        ...props.customConfig,
        [field]: rawValue,
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
                <div class="ms-custom-field-header">
                    <label class="ms-custom-label" for="ms-custom-width">Width</label>
                    <span class="ms-custom-value">{{ customConfig.cellWidth }}</span>
                </div>
                <input
                    id="ms-custom-width"
                    class="ms-custom-slider"
                    type="range"
                    :min="CUSTOM_LIMITS.minWidth"
                    :max="CUSTOM_LIMITS.maxWidth"
                    step="1"
                    :value="customConfig.cellWidth"
                    @input="updateCustomField('cellWidth', Number(($event.target as HTMLInputElement).value))"
                >
            </div>
            <div class="ms-custom-field">
                <div class="ms-custom-field-header">
                    <label class="ms-custom-label" for="ms-custom-height">Height</label>
                    <span class="ms-custom-value">{{ customConfig.cellHeight }}</span>
                </div>
                <input
                    id="ms-custom-height"
                    class="ms-custom-slider"
                    type="range"
                    :min="CUSTOM_LIMITS.minHeight"
                    :max="CUSTOM_LIMITS.maxHeight"
                    step="1"
                    :value="customConfig.cellHeight"
                    @input="updateCustomField('cellHeight', Number(($event.target as HTMLInputElement).value))"
                >
            </div>
            <div class="ms-custom-field">
                <div class="ms-custom-field-header">
                    <label class="ms-custom-label" for="ms-custom-mines">Mines</label>
                    <span class="ms-custom-value">{{ customConfig.mineCount }}</span>
                </div>
                <input
                    id="ms-custom-mines"
                    class="ms-custom-slider"
                    type="range"
                    :min="CUSTOM_LIMITS.minMines"
                    :max="maxMines"
                    step="1"
                    :value="customConfig.mineCount"
                    @input="updateCustomField('mineCount', Number(($event.target as HTMLInputElement).value))"
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
    flex-direction: column;
    gap: 8px;
    padding: 4px 2px 2px;
    width: 100%;
    max-width: 280px;
}

.ms-custom-field {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.ms-custom-field-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 8px;
}

.ms-custom-label {
    color: #000;
}

.ms-custom-value {
    color: #000;
    font-weight: bold;
    min-width: 24px;
    text-align: right;
}

.ms-custom-slider {
    width: 100%;
    height: 18px;
    margin: 0;
    accent-color: #316ac5;
    cursor: default;
}

.ms-custom-slider::-webkit-slider-runnable-track {
    height: 4px;
    background: #fff;
    border: 1px solid;
    border-color: #808080 #fff #fff #808080;
}

.ms-custom-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 11px;
    height: 18px;
    margin-top: -8px;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
}

.ms-custom-slider:active::-webkit-slider-thumb {
    border-color: #808080 #fff #fff #808080;
}

.ms-custom-slider::-moz-range-track {
    height: 4px;
    background: #fff;
    border: 1px solid;
    border-color: #808080 #fff #fff #808080;
}

.ms-custom-slider::-moz-range-thumb {
    width: 11px;
    height: 18px;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
    border-radius: 0;
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

@media (max-width: 640px) {
    .ms-settings {
        padding: 4px 8px 6px;
    }

    .ms-custom-panel {
        max-width: none;
    }

    .ms-option,
    .ms-new-game-btn {
        padding: 6px 10px;
        font-size: 12px;
    }
}
</style>
