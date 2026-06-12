<script setup lang="ts">
import type { DrawMode } from '@/lib/solitaire/types';

defineProps<{
    drawMode: DrawMode;
}>();

const emit = defineEmits<{
    'new-game': [];
    'set-draw-mode': [mode: DrawMode];
}>();
</script>

<template>
    <div class="sol-settings">
        <div class="sol-settings-row">
            <span class="sol-settings-label">Draw</span>
            <div class="sol-option-group" role="group" aria-label="Draw mode">
                <button
                    type="button"
                    class="sol-option"
                    :class="{ 'sol-option-active': drawMode === 1 }"
                    @click="emit('set-draw-mode', 1)"
                >Draw 1</button>
                <button
                    type="button"
                    class="sol-option"
                    :class="{ 'sol-option-active': drawMode === 3 }"
                    @click="emit('set-draw-mode', 3)"
                >Draw 3</button>
            </div>
        </div>

        <div class="sol-settings-actions">
            <button type="button" class="sol-new-game-btn" @click="emit('new-game')">
                New Game
            </button>
        </div>
    </div>
</template>

<style scoped>
.sol-settings {
    background: #c0c0c0;
    padding: 4px 12px 8px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.sol-settings-row {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}

.sol-settings-label {
    color: #000;
    min-width: 52px;
    flex-shrink: 0;
}

.sol-option-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.sol-option {
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

.sol-option:hover:not(.sol-option-active) {
    background: #d4d4d4;
}

.sol-option-active {
    border-color: #808080 #fff #fff #808080;
    background: #ece9d8;
    font-weight: bold;
}

.sol-settings-actions {
    display: flex;
    justify-content: center;
    padding-top: 4px;
    padding-bottom: 2px;
}

.sol-new-game-btn {
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

.sol-new-game-btn:active {
    border-color: #808080 #fff #fff #808080;
}

@media (max-width: 640px) {
    .sol-settings {
        padding: 4px 8px 6px;
    }

    .sol-option,
    .sol-new-game-btn {
        padding: 6px 10px;
        font-size: 12px;
    }
}
</style>
