<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import {
    type CellValue,
    type Grid,
    generatePuzzle,
    getConflictingCells,
    isGridComplete,
} from '@/lib/sudoku/sudoku';
import { DIFFICULTIES, GRID_SIZE, type Difficulty } from '@/lib/sudoku/types';

const props = defineProps<{
    difficulty: Difficulty;
}>();

const gameState = ref<'playing' | 'won'>('playing');
const facePressed = ref(false);
const elapsedSeconds = ref(0);
const selectedCell = ref<[number, number] | null>(null);

let timerId: ReturnType<typeof setInterval> | null = null;

const given = ref<boolean[][]>([]);
const board = ref<Grid>([]);

const faceEmoji = computed(() => {
    if (facePressed.value) return '😮';
    if (gameState.value === 'won') return '😎';
    return '🙂';
});

const conflicts = computed(() => getConflictingCells(board.value));

function formatTimer(value: number): string {
    const clamped = Math.max(0, Math.min(999, value));
    return String(clamped).padStart(3, '0');
}

function createGivenMask(grid: Grid): boolean[][] {
    return grid.map(row => row.map(value => value !== 0));
}

function startTimer() {
    if (timerId) return;
    timerId = setInterval(() => {
        if (elapsedSeconds.value < 999) elapsedSeconds.value++;
    }, 1000);
}

function stopTimer() {
    if (timerId) {
        clearInterval(timerId);
        timerId = null;
    }
}

function resetGame() {
    stopTimer();
    gameState.value = 'playing';
    facePressed.value = false;
    elapsedSeconds.value = 0;
    selectedCell.value = null;

    const { givenCount } = DIFFICULTIES[props.difficulty];
    const generated = generatePuzzle(givenCount);
    given.value = createGivenMask(generated.puzzle);
    board.value = generated.puzzle.map(row => [...row]);
    startTimer();
}

function isGiven(row: number, col: number): boolean {
    return given.value[row]?.[col] ?? false;
}

function selectCell(row: number, col: number) {
    if (gameState.value !== 'playing') return;
    selectedCell.value = [row, col];
}

function setCellValue(row: number, col: number, value: CellValue) {
    if (gameState.value !== 'playing' || isGiven(row, col)) return;

    board.value[row]![col] = value;
    checkWin();
}

function onNumberInput(value: CellValue) {
    if (!selectedCell.value) return;
    const [row, col] = selectedCell.value;
    setCellValue(row, col, value);
}

function clearSelectedCell() {
    if (!selectedCell.value) return;
    const [row, col] = selectedCell.value;
    if (isGiven(row, col)) return;
    board.value[row]![col] = 0;
}

function checkWin() {
    if (!isGridComplete(board.value)) return;
    gameState.value = 'won';
    stopTimer();
}

function cellClass(row: number, col: number): Record<string, boolean> {
    const key = `${row},${col}`;
    const isSelected =
        selectedCell.value?.[0] === row && selectedCell.value?.[1] === col;
    const sameRow = selectedCell.value?.[0] === row;
    const sameCol = selectedCell.value?.[1] === col;
    const sameBox =
        selectedCell.value !== null &&
        Math.floor(selectedCell.value[0] / 3) === Math.floor(row / 3) &&
        Math.floor(selectedCell.value[1] / 3) === Math.floor(col / 3);

    return {
        'xp-cell-given': isGiven(row, col),
        'xp-cell-selected': isSelected,
        'xp-cell-highlight': !isSelected && (sameRow || sameCol || sameBox),
        'xp-cell-conflict': conflicts.value.has(key),
        'xp-cell-box-right': col % 3 === 2 && col < GRID_SIZE - 1,
        'xp-cell-box-bottom': row % 3 === 2 && row < GRID_SIZE - 1,
    };
}

function onKeyDown(event: KeyboardEvent) {
    if (gameState.value !== 'playing') return;

    if (event.key >= '1' && event.key <= '9') {
        onNumberInput(Number(event.key) as CellValue);
        event.preventDefault();
        return;
    }
    if (event.key === 'Backspace' || event.key === 'Delete' || event.key === '0') {
        clearSelectedCell();
        event.preventDefault();
    }
}

onMounted(() => {
    resetGame();
    window.addEventListener('keydown', onKeyDown);
});

onUnmounted(() => {
    stopTimer();
    window.removeEventListener('keydown', onKeyDown);
});

watch(() => props.difficulty, () => resetGame());

defineExpose({ resetGame });
</script>

<template>
<div class="xp-client select-none" tabindex="0">
    <div class="xp-game-body">
        <div class="xp-sunken xp-status-panel">
            <div class="xp-led">{{ formatTimer(elapsedSeconds) }}</div>
            <button
                type="button"
                class="xp-face-btn"
                @click="resetGame"
                @mousedown="facePressed = true"
                @mouseup="facePressed = false"
                @mouseleave="facePressed = false"
            >{{ faceEmoji }}</button>
        </div>

        <div class="xp-board-scroll">
            <div class="xp-sunken xp-grid-panel">
                <div v-for="(row, rIndex) in board" :key="rIndex" class="flex">
                    <button
                        v-for="(value, cIndex) in row"
                        :key="cIndex"
                        type="button"
                        class="xp-cell"
                        :class="cellClass(rIndex, cIndex)"
                        @click="selectCell(rIndex, cIndex)"
                    >
                        <span
                            v-if="value !== 0"
                            class="xp-cell-value"
                        >{{ value }}</span>
                    </button>
                </div>
            </div>
        </div>

        <div class="xp-numpad">
            <button
                v-for="n in 9"
                :key="n"
                type="button"
                class="xp-num-btn"
                @click="onNumberInput(n as CellValue)"
            >{{ n }}</button>
            <button type="button" class="xp-num-btn xp-num-clear" @click="clearSelectedCell">
                Clear
            </button>
        </div>
    </div>
</div>
</template>

<style scoped>
.xp-client {
    --cell-size: 48px;
    background: #c0c0c0;
    padding: 10px 10px 6px;
    max-width: 100%;
    box-sizing: border-box;
    outline: none;
}

@media (max-width: 640px) {
    .xp-client {
        --cell-size: min(38px, calc((100vw - 3.5rem) / 9));
        width: 100%;
        padding: 8px 8px 4px;
    }

    .xp-game-body {
        width: 100%;
    }
}

.xp-game-body {
    width: fit-content;
    max-width: 100%;
    margin: auto;
}

.xp-sunken {
    border: 3px solid;
    border-color: #808080 #fff #fff #808080;
    background: #c0c0c0;
}

.xp-status-panel {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    box-sizing: border-box;
    padding: calc(var(--cell-size) * 0.2) calc(var(--cell-size) * 0.35);
    margin-bottom: calc(var(--cell-size) * 0.3);
    gap: 8px;
}

.xp-led {
    background: #000;
    color: #f00;
    font-family: 'Courier New', Courier, monospace;
    font-size: calc(var(--cell-size) * 1.1);
    font-weight: bold;
    line-height: 1;
    padding: 2px 4px;
    min-width: calc(var(--cell-size) * 2);
    text-align: right;
    letter-spacing: 1px;
}

.xp-face-btn {
    width: calc(var(--cell-size) * 1.5);
    height: calc(var(--cell-size) * 1.5);
    font-size: calc(var(--cell-size) * 0.85);
    line-height: 1;
    padding: 0;
    cursor: default;
    background: #c0c0c0;
    border: 3px solid;
    border-color: #fff #808080 #808080 #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.xp-face-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.xp-board-scroll {
    overflow: auto;
    width: 100%;
    max-width: 100%;
    max-height: calc(100vh - 12rem);
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
}

.xp-grid-panel {
    display: inline-block;
    width: max-content;
    padding: calc(var(--cell-size) * 0.2);
    border-width: 4px;
}

.xp-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    box-sizing: border-box;
    cursor: default;
    background: #fff;
    border: 1px solid #aca899;
    padding: 0;
    margin: 0;
}

.xp-cell-box-right {
    border-right-width: 3px;
    border-right-color: #000;
}

.xp-cell-box-bottom {
    border-bottom-width: 3px;
    border-bottom-color: #000;
}

.xp-cell-given .xp-cell-value {
    color: #000080;
    font-weight: bold;
}

.xp-cell-selected {
    background: #b4c8e8;
}

.xp-cell-highlight {
    background: #e8f0fc;
}

.xp-cell-conflict .xp-cell-value {
    color: #c00;
}

.xp-cell-value {
    font-size: calc(var(--cell-size) * 0.55);
    font-weight: normal;
    color: #000;
    line-height: 1;
}

.xp-numpad {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    gap: 4px;
    margin-top: calc(var(--cell-size) * 0.35);
    width: 100%;
}

.xp-num-btn {
    padding: 6px 0;
    font-size: 14px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-weight: bold;
    color: #000;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
    cursor: default;
}

.xp-num-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.xp-num-clear {
    grid-column: span 2;
    font-size: 11px;
    font-weight: normal;
}
</style>
