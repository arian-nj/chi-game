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
const highlightedCells = ref<Set<string>>(new Set())
const chordSource = ref<[number, number] | null>(null)
let suppressChordClick = false

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

function cellKey(row: number, col: number): string {
    return `${row},${col}`
}

function isCellHighlighted(row: number, col: number): boolean {
    return highlightedCells.value.has(cellKey(row, col))
}

function highlightNeighbors(row: number, col: number) {
    const next = new Set<string>()
    for (let dr = -1; dr <= 1; dr++) {
        for (let dc = -1; dc <= 1; dc++) {
            const r = row + dr
            const c = col + dc
            if (board.value[r]?.[c]) next.add(cellKey(r, c))
        }
    }
    highlightedCells.value = next
}

function clearHighlight() {
    highlightedCells.value = new Set()
    chordSource.value = null
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

function chordReveal(row: number, col: number) {
    const cell = board.value[row]?.[col]
    if (!cell?.isRevealed || cell.neighborMines === 0 || cell.isMine) return

    let neighborFlags = 0
    const hiddenNeighbors: [number, number][] = []

    for (const [dr, dc] of DIRS) {
        const r = row + dr
        const c = col + dc
        const neighbor = board.value[r]?.[c]
        if (!neighbor) continue
        if (neighbor.isFlagged) neighborFlags++
        else if (!neighbor.isRevealed) hiddenNeighbors.push([r, c])
    }

    if (neighborFlags !== cell.neighborMines) return

    for (const [r, c] of hiddenNeighbors) {
        const neighbor = board.value[r]?.[c]
        if (!neighbor) continue
        if (neighbor.isMine) {
            neighbor.isRevealed = true
            gameState.value = 'lost'
            stopTimer()
            revealAllMines()
            return
        }
        floodReveal(r, c)
    }
    checkWin()
}

function beginChord(row: number, col: number) {
    if (gameState.value !== 'playing') return

    const cell = board.value[row]?.[col]
    if (!cell?.isRevealed || cell.neighborMines === 0 || cell.isMine) return

    chordSource.value = [row, col]
    highlightNeighbors(row, col)
}

function finishChord() {
    if (!chordSource.value) return

    suppressChordClick = true
    const [row, col] = chordSource.value
    chordReveal(row, col)
    clearHighlight()
}

function onCellClick(row: number, col: number) {
    const cell = board.value[row]?.[col]

    if (cell?.isRevealed && cell.neighborMines > 0 && !cell.isMine) {
        if (suppressChordClick) {
            suppressChordClick = false
            return
        }
        highlightNeighbors(row, col)
        chordReveal(row, col)
        window.setTimeout(clearHighlight, 200)
        return
    }

    revealCell(row, col)
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
        <div
            class="xp-sunken xp-grid-panel"
            @mouseup="finishChord"
            @mouseleave="clearHighlight"
            @touchend="finishChord"
            @touchcancel="clearHighlight"
        >
            <div v-for="(row, rIndex) in board" :key="rIndex" class="flex">
                <div
                    v-for="(cell, cIndex) in row"
                    :key="cIndex"
                    class="xp-cell"
                    :class="[
                        cell.isRevealed ? 'xp-cell-revealed' : 'xp-cell-hidden',
                        { 'xp-cell-highlighted': isCellHighlighted(rIndex, cIndex) },
                    ]"
                    @click="onCellClick(rIndex, cIndex)"
                    @mousedown="beginChord(rIndex, cIndex)"
                    @touchstart="beginChord(rIndex, cIndex)"
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
    --cell-size: 26px;
    background: #c0c0c0;
    padding: 10px 10px 6px;
    max-width: 100%;
    box-sizing: border-box;
}

@media (max-width: 640px) {
    .xp-client {
        --cell-size: 34px;
        padding: 12px 12px 6px;
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
    padding: calc(var(--cell-size) * 0.35);
}

.xp-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: calc(var(--cell-size) * 0.65);
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    box-sizing: border-box;
    cursor: default;
    touch-action: manipulation;
}

.xp-cell-hidden {
    background: #c0c0c0;
    border-style: solid;
    border-width: max(3px, calc(var(--cell-size) * 0.16));
    border-color: #fff #808080 #808080 #fff;
}

.xp-cell-hidden:active {
    border-width: max(2px, calc(var(--cell-size) * 0.1));
    border-color: #808080;
}

.xp-cell-revealed {
    background: #c0c0c0;
    border-style: solid;
    border-width: max(2px, calc(var(--cell-size) * 0.08));
    border-color: #808080;
}

.xp-cell-highlighted {
    animation: xp-chord-wiggle 0.12s ease-in-out infinite alternate;
}

.xp-cell-highlighted.xp-cell-hidden {
    border-width: max(2px, calc(var(--cell-size) * 0.1));
    border-color: #808080;
}

.xp-cell-highlighted.xp-cell-revealed {
    background: #a8a8a8;
    box-shadow: inset 0 0 0 1px #ffd700;
}

@keyframes xp-chord-wiggle {
    from {
        transform: translateX(-1px);
    }
    to {
        transform: translateX(1px);
    }
}

.xp-flag {
    font-size: calc(var(--cell-size) * 0.55);
    line-height: 1;
}
</style>
