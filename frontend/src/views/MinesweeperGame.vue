<script setup lang="ts">
import { ref } from 'vue';


const props = defineProps<{
    cellWidth: number;
    cellHeight: number;
    mineCount: number;
}>()

interface Cell {
    isMine: boolean;
    neighbortMines: number;

    isFlagged: boolean;
    isRevealed: boolean;
}

const board = ref<Cell[][]>([])

board.value = Array.from({length:props.cellHeight}, ()=>Array.from({length:props.cellWidth}, ()=>({isMine:false, neighbortMines:0, isFlagged:false, isRevealed:false})))

function placeBombsRandomly() {
    const emptyCells: { row: number; col: number }[] = [];
    for (let row = 0; row < props.cellHeight; row++) {
        for (let col = 0; col < props.cellWidth; col++) {
            emptyCells.push({ row, col });
        }
    }

    for (let i = 0; i < props.mineCount && emptyCells.length > 0; i++) {
        const index = Math.floor(Math.random() * emptyCells.length);
        const { row, col } = emptyCells.splice(index, 1)[0]!;
        board.value[row]![col]!.isMine = true;
    }
}

function calculateNeighbortMines() {
    for (let row = 0; row < props.cellHeight; row++) {
        for (let col = 0; col < props.cellWidth; col++) {
            const cell = board.value[row]?.[col];
            if (cell?.isMine) {
                for (let r = row - 1; r <= row + 1; r++) {
                    for (let c = col - 1; c <= col + 1; c++) {
                        const neighbor = board.value[r]?.[c];
                        if (neighbor && !neighbor.isMine) {
                            neighbor.neighbortMines++;
                        }
                    }
                }
            }
        }
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


// const startTime = performance.now()
placeBombsRandomly()
calculateNeighbortMines()
// const endTime = performance.now()
// const timeTaken = endTime - startTime
// console.log(`Time taken to place bombs and calculate neighbort mines: ${timeTaken} milliseconds`)

function revealCell(row: number, col: number) {
    const cell = board.value[row]?.[col];
    if (cell) {
        cell.isRevealed = true;
    }
}

function rightClickCell(row: number, col: number) {
    const cell = board.value[row]?.[col];
    if (cell) {
        cell.isFlagged = !cell.isFlagged;
    }
}

</script>

<template>

<div class="bg-gray-50 w-full h-full">
    <div v-for="row,rIndex in board" :key="rIndex" class="flex">
        <div v-for="cell,hIndex in row" :key="hIndex">
            <div class="bg-[#c2c2c2] w-12 h-12 flex items-center justify-center text-2xl font-bold
                border border-[#828282]
            ">
                <div @click="revealCell(rIndex, hIndex)" @contextmenu.prevent="rightClickCell(rIndex, hIndex)" v-if="!cell.isRevealed" class="w-full h-full flex items-center justify-center border-4 border-gray-300">{{ cell.isFlagged ? '🚩' : '' }}</div>
                <span v-else-if="cell.isRevealed && cell.isMine">💣</span>
                <span v-else-if="cell.isRevealed" :class="getNumbersColor(cell.neighbortMines)">{{ cell.neighbortMines == 0 ? '' : cell.neighbortMines }}</span>
            </div>
        </div>
    </div>

    <!-- <div class="bg-gray-400 w-10 h-10"></div> -->
</div>

</template>

<style scoped>

</style>