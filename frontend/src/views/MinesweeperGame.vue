<script setup lang="ts">
import { ref } from 'vue';
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

interface Cell {
    isMine: boolean;
    neighborMines: number;

    isFlagged: boolean;
    isRevealed: boolean;
}

const board = ref<Cell[][]>([])

board.value = Array.from({length:props.cellHeight}, ()=>Array.from({length:props.cellWidth}, ()=>({isMine:false, neighborMines:0, isFlagged:false, isRevealed:false})))

const DIRS: [number, number][] = [
  [-1, -1], [-1, 0], [-1, 1],
  [0, -1],           [0, 1],
  [1, -1],  [1, 0],  [1, 1],
]

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
        case 1:
            return 'text-[#110cbd]';
        case 2:
            return 'text-[#155f10]';
        case 3:
            return 'text-[#ff0000]';
        case 4:
            return 'text-[#00007b]';
        case 5:
            return 'text-[#ffa500]';
        case 6:
            return 'text-[#00ffff]';
        case 7:
            return 'text-[#800080]';
        case 8:
            return 'text-[#000000]';
        default:
            return 'text-[#000000]';
    }
}

function revealCell(row: number, col: number) {
    if (gameState.value !== 'playing') return;

    const cell = board.value[row]?.[col];

    if (!cell || cell.isRevealed || cell.isFlagged) return;

    if (!minesPlaced.value) {
        placeMinesUntilSolvable(row, col)
        minesPlaced.value = true
    }
    
    if (cell.isMine) {
        cell.isRevealed = true
        gameState.value = 'lost'
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
    // Only expand from empty cells
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
  if (allSafeRevealed) gameState.value = 'won'
}
function revealAllMines() {
  for (const row of board.value) {
    for (const cell of row) {
      if (cell.isMine) cell.isRevealed = true
    }
  }
}

</script>

<template>

<div class="bg-gray-50 w-full h-full">
    <div v-for="row,rIndex in board" :key="rIndex" class="flex">
        <div v-for="cell,hIndex in row" :key="hIndex">
            <div class="bg-[#c2c2c2] w-8 h-8 flex items-center justify-center text-2xl font-bold
                border border-[#828282]
            ">
                <div @click="revealCell(rIndex, hIndex)" @contextmenu.prevent="rightClickCell(rIndex, hIndex)" v-if="!cell.isRevealed" 
                    :class='["w-full h-full flex items-center justify-center border-4 border-t-[#ffffff] border-l-[#ffffff]  border-b-[#7b7b7b] border-r-[#7b7b7b]",
                        cell.isFlagged?"text-lg":"text-2xl"]'>
                        {{ cell.isFlagged ? '🚩' : '' }}
                </div>
                <span v-else-if="cell.isRevealed && cell.isMine">💣</span>
                <span v-else-if="cell.isRevealed" :class="getNumbersColor(cell.neighborMines)">{{ cell.neighborMines == 0 ? '' : cell.neighborMines }}</span>
            </div>
        </div>
    </div>

    <!-- <div class="bg-gray-400 w-10 h-10"></div> -->
</div>

</template>

<style scoped>

</style>