<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
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
let mineRevealTimeoutIds: ReturnType<typeof setTimeout>[] = []

interface Cell {
    isMine: boolean;
    neighborMines: number;
    isFlagged: boolean;
    isRevealed: boolean;
}

const board = ref<Cell[][]>(createEmptyBoard())
const highlightedCells = ref<Set<string>>(new Set())
const previewSource = ref<[number, number] | null>(null)
let suppressNumberClick = false

const MIN_CELL_ZOOM = 0.25
const MAX_CELL_ZOOM = 2
const CELL_ZOOM_STEP = 0.25
const SMALL_SCREEN_QUERY = '(max-width: 640px)'

const baseCellSize = ref(26)
const cellZoom = ref(1)
const isPinching = ref(false)
const boardScrollRef = ref<HTMLElement | null>(null)
let pinchStartDistance = 0
let pinchStartZoom = 1
let smallScreenQuery: MediaQueryList | null = null

const BOARD_PAN_THRESHOLD = 5
let isBoardPanning = false
let suppressBoardClick = false
let panStartX = 0
let panStartY = 0
let panScrollLeft = 0
let panScrollTop = 0

const effectiveCellSize = computed(() => Math.round(baseCellSize.value * cellZoom.value))
const zoomPercent = computed(() => Math.round(cellZoom.value * 100))
const boardStyle = computed(() => ({
    '--cell-size': `${effectiveCellSize.value}px`,
}))
const canZoomOut = computed(() => cellZoom.value > MIN_CELL_ZOOM + 0.001)
const canZoomIn = computed(() => cellZoom.value < MAX_CELL_ZOOM - 0.001)

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

function stopMineRevealAnimation() {
    for (const id of mineRevealTimeoutIds) clearTimeout(id)
    mineRevealTimeoutIds = []
}

function resetGame() {
    stopTimer()
    stopMineRevealAnimation()
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

function clampZoom(value: number): number {
    return Math.min(MAX_CELL_ZOOM, Math.max(MIN_CELL_ZOOM, value))
}

function updateBaseCellSize() {
    if (!window.matchMedia(SMALL_SCREEN_QUERY).matches) {
        baseCellSize.value = 26
        return
    }

    const horizontalPadding = 56
    const boardChrome = 14
    const available = window.innerWidth - horizontalPadding
    const fitSize = Math.floor((available - boardChrome) / props.cellWidth)
    baseCellSize.value = Math.min(48, Math.max(22, fitSize))
}

function zoomIn() {
    cellZoom.value = clampZoom(
        Math.round((cellZoom.value + CELL_ZOOM_STEP) * 100) / 100,
    )
}

function zoomOut() {
    cellZoom.value = clampZoom(
        Math.round((cellZoom.value - CELL_ZOOM_STEP) * 100) / 100,
    )
}

function touchDistance(touches: TouchList): number {
    const first = touches[0]
    const second = touches[1]
    if (!first || !second) return 0
    return Math.hypot(first.clientX - second.clientX, first.clientY - second.clientY)
}

function onBoardTouchStart(event: TouchEvent) {
    if (event.touches.length !== 2) return

    isPinching.value = true
    pinchStartDistance = touchDistance(event.touches)
    pinchStartZoom = cellZoom.value
    clearHighlight()
}

function onBoardTouchMove(event: TouchEvent) {
    if (!isPinching.value || event.touches.length < 2 || pinchStartDistance <= 0) return

    event.preventDefault()
    const distance = touchDistance(event.touches)
    cellZoom.value = clampZoom(pinchStartZoom * (distance / pinchStartDistance))
}

function onBoardTouchEnd(event: TouchEvent) {
    if (event.touches.length < 2) {
        isPinching.value = false
        pinchStartDistance = 0
    }
}

function onBoardMouseDown(event: MouseEvent) {
    if (event.button !== 0) return
    if (isPinching.value) return

    const boardScroll = boardScrollRef.value
    if (!boardScroll) return

    isBoardPanning = false
    panStartX = event.clientX
    panStartY = event.clientY
    panScrollLeft = boardScroll.scrollLeft
    panScrollTop = boardScroll.scrollTop

    document.addEventListener('mousemove', onBoardMouseMove)
    document.addEventListener('mouseup', onBoardMouseUp)
}

function onBoardMouseMove(event: MouseEvent) {
    const boardScroll = boardScrollRef.value
    if (!boardScroll) return

    const dx = event.clientX - panStartX
    const dy = event.clientY - panStartY

    if (!isBoardPanning) {
        if (Math.abs(dx) < BOARD_PAN_THRESHOLD && Math.abs(dy) < BOARD_PAN_THRESHOLD) return

        isBoardPanning = true
        clearHighlight()
        boardScroll.classList.add('xp-board-scroll-dragging')
    }

    event.preventDefault()
    boardScroll.scrollLeft = panScrollLeft - dx
    boardScroll.scrollTop = panScrollTop - dy
}

function stopBoardPan() {
    document.removeEventListener('mousemove', onBoardMouseMove)
    document.removeEventListener('mouseup', onBoardMouseUp)

    boardScrollRef.value?.classList.remove('xp-board-scroll-dragging')
}

function onBoardMouseUp() {
    if (isBoardPanning) {
        suppressBoardClick = true
    }

    isBoardPanning = false
    stopBoardPan()
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
    previewSource.value = null
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

function beginNeighborPreview(row: number, col: number, event?: MouseEvent | TouchEvent) {
    if (event && 'touches' in event && event.touches.length > 1) return
    if (isPinching.value) return
    if (gameState.value !== 'playing') return

    const cell = board.value[row]?.[col]
    if (!cell?.isRevealed || cell.neighborMines === 0 || cell.isMine) return

    previewSource.value = [row, col]
    highlightNeighbors(row, col)
}

function finishNeighborPreview() {
    if (isPinching.value || isBoardPanning) return
    if (!previewSource.value) return

    suppressNumberClick = true
    clearHighlight()
}

function onCellClick(row: number, col: number) {
    if (suppressBoardClick) {
        suppressBoardClick = false
        return
    }

    const cell = board.value[row]?.[col]

    if (cell?.isRevealed && cell.neighborMines > 0 && !cell.isMine) {
        if (suppressNumberClick) {
            suppressNumberClick = false
            return
        }
        highlightNeighbors(row, col)
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
    celebrateWin()
  }
}

function celebrateWin() {
    gameState.value = 'won'
    stopTimer()

    for (const row of board.value) {
        for (const cell of row) {
            if (cell.isMine) cell.isFlagged = true
        }
    }
}

function revealAllMines() {
    stopMineRevealAnimation()

    const hiddenMines: [number, number][] = []
    for (let row = 0; row < props.cellHeight; row++) {
        for (let col = 0; col < props.cellWidth; col++) {
            const cell = board.value[row]![col]!
            if (cell.isMine && !cell.isRevealed) hiddenMines.push([row, col])
        }
    }

    if (hiddenMines.length === 0) return

    const stepMs = Math.min(100, Math.max(40, Math.floor(2500 / hiddenMines.length)))

    hiddenMines.forEach(([row, col], index) => {
        const id = window.setTimeout(() => {
            if (gameState.value !== 'lost') return
            const cell = board.value[row]?.[col]
            if (cell?.isMine) cell.isRevealed = true
        }, index * stepMs)
        mineRevealTimeoutIds.push(id)
    })
}

watch(
    () => [props.cellWidth, props.cellHeight, props.mineCount],
    () => {
        resetGame()
        updateBaseCellSize()
    },
)

function onViewportChange() {
    updateBaseCellSize()
}

onMounted(() => {
    updateBaseCellSize()
    smallScreenQuery = window.matchMedia(SMALL_SCREEN_QUERY)
    smallScreenQuery.addEventListener('change', onViewportChange)
    window.addEventListener('resize', onViewportChange)

    const boardScroll = boardScrollRef.value
    if (boardScroll) {
        boardScroll.addEventListener('touchstart', onBoardTouchStart, { passive: true })
        boardScroll.addEventListener('touchmove', onBoardTouchMove, { passive: false })
        boardScroll.addEventListener('touchend', onBoardTouchEnd, { passive: true })
        boardScroll.addEventListener('touchcancel', onBoardTouchEnd, { passive: true })
        boardScroll.addEventListener('mousedown', onBoardMouseDown)
    }
})

onUnmounted(() => {
    stopTimer()
    stopMineRevealAnimation()
    smallScreenQuery?.removeEventListener('change', onViewportChange)
    window.removeEventListener('resize', onViewportChange)

    stopBoardPan()

    const boardScroll = boardScrollRef.value
    if (boardScroll) {
        boardScroll.removeEventListener('touchstart', onBoardTouchStart)
        boardScroll.removeEventListener('touchmove', onBoardTouchMove)
        boardScroll.removeEventListener('touchend', onBoardTouchEnd)
        boardScroll.removeEventListener('touchcancel', onBoardTouchEnd)
        boardScroll.removeEventListener('mousedown', onBoardMouseDown)
    }
})

defineExpose({ resetGame })
</script>

<template>
<div class="xp-client select-none">
    <div class="xp-game-body" :class="{ 'xp-game-won': gameState === 'won' }">
        <div class="xp-sunken xp-status-panel">
            <div class="xp-led">{{ formatCounter(mineCounter) }}</div>
            <button
                type="button"
                class="xp-face-btn"
                :class="{ 'xp-face-btn-won': gameState === 'won' }"
                @click="resetGame"
                @mousedown="facePressed = true"
                @mouseup="facePressed = false"
                @mouseleave="facePressed = false"
            >{{ faceEmoji }}</button>
            <div class="xp-led">{{ formatCounter(elapsedSeconds) }}</div>
        </div>

        <div class="xp-zoom-controls" aria-label="Board zoom controls">
            <button
                type="button"
                class="xp-zoom-btn"
                :disabled="!canZoomOut"
                aria-label="Zoom out"
                @click="zoomOut"
            >−</button>
            <span class="xp-zoom-label">{{ zoomPercent }}%</span>
            <button
                type="button"
                class="xp-zoom-btn"
                :disabled="!canZoomIn"
                aria-label="Zoom in"
                @click="zoomIn"
            >+</button>
        </div>

        <div class="xp-board-area">
            <div
                v-if="gameState === 'won'"
                class="xp-win-banner"
                role="status"
                aria-live="polite"
            >
                <span class="xp-win-title">You Win!</span>
                <span class="xp-win-time">Cleared in {{ formatCounter(elapsedSeconds) }}s</span>
            </div>

            <div ref="boardScrollRef" class="xp-board-scroll" :style="boardStyle">
            <div
                class="xp-sunken xp-grid-panel"
                :class="{ 'xp-grid-panel-won': gameState === 'won' }"
                @mouseup="finishNeighborPreview"
                @mouseleave="clearHighlight"
                @touchend="finishNeighborPreview"
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
                            { 'xp-cell-mine-reveal': cell.isRevealed && cell.isMine },
                        ]"
                        @click="onCellClick(rIndex, cIndex)"
                        @mousedown="beginNeighborPreview(rIndex, cIndex)"
                        @touchstart="beginNeighborPreview(rIndex, cIndex, $event)"
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
</div>
</template>

<style scoped>
.xp-client {
    --ui-cell-size: 26px;
    background: #c0c0c0;
    padding: 10px 10px 6px;
    max-width: 100%;
    box-sizing: border-box;
}

@media (max-width: 640px) {
    .xp-client {
        --ui-cell-size: 34px;
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
    padding: calc(var(--ui-cell-size) * 0.2) calc(var(--ui-cell-size) * 0.35);
    margin-bottom: calc(var(--ui-cell-size) * 0.3);
    transition: box-shadow 0.3s ease;
}

.xp-game-won .xp-status-panel {
    box-shadow: 0 0 0 2px #2e8b2e, 0 0 12px rgba(46, 139, 46, 0.55);
}

.xp-led {
    background: #000;
    color: #f00;
    font-family: 'Courier New', Courier, monospace;
    font-size: calc(var(--ui-cell-size) * 1.1);
    font-weight: bold;
    line-height: 1;
    padding: 2px 4px;
    min-width: calc(var(--ui-cell-size) * 2);
    text-align: right;
    letter-spacing: 1px;
}

.xp-face-btn {
    width: calc(var(--ui-cell-size) * 1.5);
    height: calc(var(--ui-cell-size) * 1.5);
    font-size: calc(var(--ui-cell-size) * 0.85);
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

.xp-face-btn-won {
    background: #fff9b0;
    animation: xp-win-face 0.55s ease-in-out infinite alternate;
}

@keyframes xp-win-face {
    from {
        transform: scale(1);
    }
    to {
        transform: scale(1.12);
    }
}

.xp-zoom-controls {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: calc(var(--ui-cell-size) * 0.35);
    margin-bottom: calc(var(--ui-cell-size) * 0.25);
}

.xp-zoom-btn {
    width: calc(var(--ui-cell-size) * 1.15);
    height: calc(var(--ui-cell-size) * 1.15);
    font-size: calc(var(--ui-cell-size) * 0.75);
    line-height: 1;
    padding: 0;
    background: #c0c0c0;
    border: 3px solid;
    border-color: #fff #808080 #808080 #fff;
    color: #000;
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
}

.xp-zoom-btn:disabled {
    color: #808080;
}

.xp-zoom-btn:not(:disabled):active {
    border-color: #808080 #fff #fff #808080;
}

.xp-zoom-label {
    min-width: calc(var(--ui-cell-size) * 1.8);
    text-align: center;
    font-size: calc(var(--ui-cell-size) * 0.45);
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    color: #000;
}


.xp-board-area {
    position: relative;
}

.xp-win-banner {
    position: absolute;
    inset: 0;
    z-index: 5;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    pointer-events: none;
    background: rgba(0, 80, 0, 0.28);
    animation: xp-win-banner-in 0.45s ease-out;
}

.xp-win-title {
    padding: 8px 20px;
    background: linear-gradient(180deg, #fff9b0 0%, #ffd700 55%, #e6b800 100%);
    border: 3px solid;
    border-color: #fff #8b7500 #8b7500 #fff;
    color: #1a4d1a;
    font-size: clamp(1.25rem, 4vw, 1.75rem);
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    text-shadow: 1px 1px 0 #fff;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
    animation: xp-win-title-pop 0.5s ease-out;
}

.xp-win-time {
    padding: 4px 12px;
    background: rgba(0, 0, 0, 0.72);
    color: #7fff7f;
    font-family: 'Courier New', Courier, monospace;
    font-size: clamp(0.85rem, 2.5vw, 1rem);
    font-weight: bold;
    letter-spacing: 1px;
}

@keyframes xp-win-banner-in {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

@keyframes xp-win-title-pop {
    0% {
        transform: scale(0.5);
        opacity: 0;
    }
    70% {
        transform: scale(1.08);
        opacity: 1;
    }
    100% {
        transform: scale(1);
    }
}

.xp-board-scroll {
    overflow: auto;
    width: 100%;
    max-width: 100%;
    max-height: calc(100vh - 7rem);
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
    touch-action: pan-x pan-y;
}

.xp-board-scroll-dragging {
    cursor: grabbing;
    user-select: none;
}

.xp-grid-panel {
    display: inline-block;
    width: max-content;
    padding: calc(var(--cell-size) * 0.35);
    transition: box-shadow 0.3s ease;
}

.xp-grid-panel-won {
    box-shadow: 0 0 0 3px #2e8b2e, 0 0 18px rgba(80, 220, 80, 0.65);
    animation: xp-win-grid-glow 1s ease-in-out infinite alternate;
}

@keyframes xp-win-grid-glow {
    from {
        box-shadow: 0 0 0 3px #2e8b2e, 0 0 12px rgba(80, 220, 80, 0.45);
    }
    to {
        box-shadow: 0 0 0 3px #4caf50, 0 0 24px rgba(120, 255, 120, 0.75);
    }
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

.xp-cell-mine-reveal {
    animation: xp-mine-pop 0.18s ease-out;
}

@keyframes xp-mine-pop {
    from {
        transform: scale(0.4);
    }
    to {
        transform: scale(1);
    }
}

.xp-cell-highlighted {
    animation: xp-chord-wiggle 0.12s ease-in-out infinite alternate;
}

.xp-cell-highlighted.xp-cell-hidden {
    border-width: max(2px, calc(var(--cell-size) * 0.1));
    border-color: #458b19;
}

.xp-cell-highlighted.xp-cell-revealed {
    background: #9fdfae47;
    box-shadow: inset 0 0 0 2px #2e8b2e;
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
