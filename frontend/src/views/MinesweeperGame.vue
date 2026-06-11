<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue';
import {
    computeNeighborMines,
    generateSolvableMines,
} from '@/lib/minesweeper/minesweeperSolver'

const props = defineProps<{
    cellWidth: number;
    cellHeight: number;
    mineCount: number;
}>()

const gameState = ref<'playing' | 'won' | 'lost'>('playing');
const minesPlaced = ref(false);
const facePressed = ref(false);
const elapsedSeconds = ref(0);
let timerId: ReturnType<typeof setInterval> | null = null;

interface Cell {
    isMine: boolean;
    neighborMines: number;
    isFlagged: boolean;
    isRevealed: boolean;
}

const board = ref<Cell[][]>(createEmptyBoard())

function createEmptyBoard(): Cell[][] {
    return Array.from({ length: props.cellHeight }, () =>
        Array.from({ length: props.cellWidth }, () => ({
            isMine: false,
            neighborMines: 0,
            isFlagged: false,
            isRevealed: false,
        }))
    )
}

const DIRS: [number, number][] = [
  [-1, -1], [-1, 0], [-1, 1],
  [0, -1],           [0, 1],
  [1, -1],  [1, 0],  [1, 1],
]

const flagCount = computed(() =>
    board.value.reduce((sum, row) => sum + row.filter(c => c.isFlagged).length, 0)
)

const mineCounter = computed(() => {
    const remaining = props.mineCount - flagCount.value
    return Math.max(-99, Math.min(999, remaining))
})

const faceEmoji = computed(() => {
    if (facePressed.value) return '😮'
    if (gameState.value === 'won') return '😎'
    if (gameState.value === 'lost') return '😵'
    return '🙂'
})

function formatCounter(value: number): string {
    const clamped = Math.max(-99, Math.min(999, value))
    const sign = clamped < 0 ? '-' : ''
    return sign + String(Math.abs(clamped)).padStart(3, '0')
}

function startTimer() {
    if (timerId) return
    timerId = setInterval(() => {
        if (elapsedSeconds.value < 999) elapsedSeconds.value++
    }, 1000)
}

function stopTimer() {
    if (timerId) {
        clearInterval(timerId)
        timerId = null
    }
}

function resetGame() {
    stopTimer()
    gameState.value = 'playing'
    minesPlaced.value = false
    facePressed.value = false
    elapsedSeconds.value = 0
    board.value = createEmptyBoard()
}

function applyGeneratedMines(mines: boolean[][]) {
  const neighborCounts = computeNeighborMines(mines, props.cellHeight, props.cellWidth)

  for (let row = 0; row < props.cellHeight; row++) {
    for (let col = 0; col < props.cellWidth; col++) {
      const cell = board.value[row]![col]!
      cell.isMine = mines[row]![col]!
      cell.neighborMines = neighborCounts[row]![col]!
    }
  }
}

function placeMinesUntilSolvable(safeRow: number, safeCol: number) {
  const { mines, solvable } = generateSolvableMines(
    props.cellHeight,
    props.cellWidth,
    props.mineCount,
    safeRow,
    safeCol,
  )

  applyGeneratedMines(mines)

  if (!solvable) {
    console.warn('Could not generate a logic-solvable board within the attempt limit')
  }
}

function getNumbersColor(number: number) {
    switch (number) {
        case 1: return 'text-[#0000ff]';
        case 2: return 'text-[#008000]';
        case 3: return 'text-[#ff0000]';
        case 4: return 'text-[#000080]';
        case 5: return 'text-[#800000]';
        case 6: return 'text-[#008080]';
        case 7: return 'text-[#000000]';
        case 8: return 'text-[#808080]';
        default: return 'text-[#000000]';
    }
}

function revealCell(row: number, col: number) {
    if (gameState.value !== 'playing') return;

    const cell = board.value[row]?.[col];

    if (!cell || cell.isRevealed || cell.isFlagged) return;

    if (!minesPlaced.value) {
        placeMinesUntilSolvable(row, col)
        minesPlaced.value = true
        startTimer()
    }
    
    if (cell.isMine) {
        cell.isRevealed = true
        gameState.value = 'lost'
        stopTimer()
        revealAllMines()
        return
    }

    floodReveal(row, col)
    checkWin()
}

function floodReveal(startRow: number, startCol: number) {
  const queue: [number, number][] = [[startRow, startCol]]
  while (queue.length > 0) {
    const [row, col] = queue.shift()!
    const cell = board.value[row]?.[col]
    if (!cell || cell.isRevealed || cell.isFlagged) continue
    cell.isRevealed = true
    if (cell.neighborMines !== 0) continue
    for (const [dr, dc] of DIRS) {
      const r = row + dr
      const c = col + dc
      const neighbor = board.value[r]?.[c]
      if (neighbor && !neighbor.isRevealed && !neighbor.isFlagged) {
        queue.push([r, c])
      }
    }
  }
}

function rightClickCell(row: number, col: number) {
  if (gameState.value !== 'playing') return

  const cell = board.value[row]?.[col]
  if (!cell || cell.isRevealed) return

  cell.isFlagged = !cell.isFlagged
}

function checkWin() {
  const allSafeRevealed = board.value.every(row =>
    row.every(cell => cell.isMine || cell.isRevealed)
  )
  if (allSafeRevealed) {
    gameState.value = 'won'
    stopTimer()
  }
}

function revealAllMines() {
  for (const row of board.value) {
    for (const cell of row) {
      if (cell.isMine) cell.isRevealed = true
    }
  }
}

watch(
    () => [props.cellWidth, props.cellHeight, props.mineCount],
    () => resetGame(),
)

onUnmounted(stopTimer)

defineExpose({ resetGame })
</script>

<template>
<div class="xp-client select-none">
    <div class="xp-game-body">
    <div class="xp-sunken xp-status-panel">
        <div class="xp-led">{{ formatCounter(mineCounter) }}</div>
        <button
            type="button"
            class="xp-face-btn"
            @click="resetGame"
            @mousedown="facePressed = true"
            @mouseup="facePressed = false"
            @mouseleave="facePressed = false"
        >{{ faceEmoji }}</button>
        <div class="xp-led">{{ formatCounter(elapsedSeconds) }}</div>
    </div>

    <div class="xp-board-scroll">
        <div class="xp-sunken xp-grid-panel">
            <div v-for="(row, rIndex) in board" :key="rIndex" class="flex">
                <div
                    v-for="(cell, cIndex) in row"
                    :key="cIndex"
                    class="xp-cell"
                    :class="cell.isRevealed ? 'xp-cell-revealed' : 'xp-cell-hidden'"
                    @click="revealCell(rIndex, cIndex)"
                    @contextmenu.prevent="rightClickCell(rIndex, cIndex)"
                >
                    <template v-if="!cell.isRevealed">
                        <span v-if="cell.isFlagged" class="xp-flag">🚩</span>
                    </template>
                    <template v-else-if="cell.isMine">💣</template>
                    <span
                        v-else-if="cell.neighborMines > 0"
                        :class="getNumbersColor(cell.neighborMines)"
                    >{{ cell.neighborMines }}</span>
                </div>
            </div>
        </div>
    </div>
    </div>
</div>
</template>

<style scoped>
.xp-client {
    --cell-size: 20px;
    background: #c0c0c0;
    padding: 8px;
    max-width: 100%;
    box-sizing: border-box;
}

.xp-game-body {
    width: fit-content;
    max-width: 100%;
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
    padding: 5px 9px;
    margin-bottom: 8px;
}

.xp-led {
    background: #000;
    color: #f00;
    font-family: 'Courier New', Courier, monospace;
    font-size: 28px;
    font-weight: bold;
    line-height: 1;
    padding: 2px 4px;
    min-width: 50px;
    text-align: right;
    letter-spacing: 1px;
}

.xp-face-btn {
    width: 32px;
    height: 32px;
    font-size: 20px;
    line-height: 1;
    padding: 0;
    cursor: default;
    background: #c0c0c0;
    border: 3px solid;
    border-color: #fff #808080 #808080 #fff;
    display: flex;
    align-items: center;
    justify-content: center;
}

.xp-face-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.xp-board-scroll {
    overflow: auto;
    width: 100%;
    max-width: 100%;
    max-height: calc(100vh - 7rem);
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
}

.xp-grid-panel {
    display: inline-block;
    width: max-content;
    padding: 8px;
}

.xp-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    box-sizing: border-box;
    cursor: default;
}

.xp-cell-hidden {
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
}

.xp-cell-hidden:active {
    border-width: 1px;
    border-color: #808080;
}

.xp-cell-revealed {
    background: #c0c0c0;
    border: 1px solid #808080;
}

.xp-flag {
    font-size: 12px;
    line-height: 1;
}
</style>
