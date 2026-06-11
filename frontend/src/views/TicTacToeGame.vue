<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import {
    checkWinner,
    createEmptyBoard,
    getBotMove,
    isDraw,
    type Cell,
    type Player,
} from '@/lib/tictactoe/tictactoe';
import type { BoardSize, BotDifficulty } from '@/lib/tictactoe/types';

const props = defineProps<{
    boardSize: BoardSize;
    isBot: boolean;
    botDifficulty: BotDifficulty;
}>();

const HUMAN: Player = 'X';
const BOT: Player = 'O';

const board = ref<Cell[]>(createEmptyBoard(props.boardSize));
const currentPlayer = ref<Player>(HUMAN);
const winningLine = ref<number[]>([]);
const isBotThinking = ref(false);
const facePressed = ref(false);
const moveCount = ref(0);

const result = computed(() => checkWinner(board.value, props.boardSize));
const hasDraw = computed(() => isDraw(board.value, props.boardSize));
const isGameOver = computed(() => result.value !== null || hasDraw.value);

const turnIndicator = computed(() => {
    if (result.value) return result.value.winner;
    if (hasDraw.value) return '—';
    if (isBotThinking.value) return 'O';
    return currentPlayer.value;
});

function formatCounter(value: number): string {
    return String(Math.min(999, value)).padStart(3, '0');
}

const faceEmoji = computed(() => {
    if (facePressed.value) return '😮';
    if (result.value) {
        if (props.isBot) {
            return result.value.winner === HUMAN ? '😎' : '😵';
        }
        return '😎';
    }
    if (hasDraw.value) return '😐';
    if (isBotThinking.value) return '🤔';
    return '🙂';
});

function resetGame() {
    board.value = createEmptyBoard(props.boardSize);
    currentPlayer.value = HUMAN;
    winningLine.value = [];
    isBotThinking.value = false;
    facePressed.value = false;
    moveCount.value = 0;
}

function isWinningCell(index: number): boolean {
    return winningLine.value.includes(index);
}

function placeMove(index: number, player: Player) {
    if (board.value[index] !== null || isGameOver.value) return;

    const nextBoard = [...board.value];
    nextBoard[index] = player;
    board.value = nextBoard;
    moveCount.value++;

    const gameResult = checkWinner(board.value, props.boardSize);
    if (gameResult) {
        winningLine.value = gameResult.line;
        return;
    }

    if (isDraw(board.value, props.boardSize)) return;

    currentPlayer.value = player === HUMAN ? BOT : HUMAN;
}

function handleCellClick(index: number) {
    if (isGameOver.value || isBotThinking.value) return;
    if (props.isBot && currentPlayer.value !== HUMAN) return;
    if (board.value[index] !== null) return;

    placeMove(index, currentPlayer.value);
}

function runBotMove() {
    if (!props.isBot || isGameOver.value || currentPlayer.value !== BOT) return;

    isBotThinking.value = true;
    window.setTimeout(() => {
        const move = getBotMove(
            [...board.value],
            props.boardSize,
            BOT,
            HUMAN,
            props.botDifficulty,
        );
        isBotThinking.value = false;
        if (move >= 0) placeMove(move, BOT);
    }, 350);
}

watch(
    () => [currentPlayer.value, isGameOver.value, props.isBot] as const,
    () => runBotMove(),
    { flush: 'post' },
);

watch(
    () => [props.boardSize, props.isBot, props.botDifficulty],
    () => resetGame(),
);

defineExpose({ resetGame });
</script>

<template>
    <div
        class="xp-client select-none"
        :class="boardSize === 5 ? 'xp-board-5' : 'xp-board-3'"
        :style="{ '--grid-cols': boardSize }"
    >
        <div class="xp-game-body">
            <div class="xp-sunken xp-status-panel">
                <div class="xp-led">{{ turnIndicator }}</div>
                <button
                    type="button"
                    class="xp-face-btn"
                    @click="resetGame"
                    @mousedown="facePressed = true"
                    @mouseup="facePressed = false"
                    @mouseleave="facePressed = false"
                >{{ faceEmoji }}</button>
                <div class="xp-led">{{ formatCounter(moveCount) }}</div>
            </div>

            <div class="xp-sunken xp-grid-panel">
                <button
                    v-for="(cell, index) in board"
                    :key="index"
                    type="button"
                    class="xp-cell"
                    :class="{
                        'xp-cell-empty': cell === null,
                        'xp-cell-filled': cell !== null,
                        'xp-cell-winning': isWinningCell(index),
                        'xp-cell-x': cell === 'X',
                        'xp-cell-o': cell === 'O',
                    }"
                    :disabled="
                        cell !== null
                            || isGameOver
                            || isBotThinking
                            || (isBot && currentPlayer !== HUMAN)
                    "
                    @click="handleCellClick(index)"
                >
                    {{ cell ?? '' }}
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.xp-client {
    --cell-size: 64px;
    --panel-pad: calc(var(--cell-size) * 0.35);
    --board-width: calc(var(--cell-size) * var(--grid-cols) + var(--panel-pad) * 2 + 6px);
    width: fit-content;
    background: #c0c0c0;
    padding: 10px;
    box-sizing: border-box;
}

.xp-board-5 {
    --cell-size: 44px;
}

@media (max-width: 640px) {
    .xp-client {
        --cell-size: 72px;
        padding: 12px;
    }

    .xp-board-5 {
        --cell-size: 48px;
    }
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
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: var(--board-width);
    box-sizing: border-box;
    padding: calc(var(--cell-size) * 0.15) calc(var(--cell-size) * 0.2);
    margin-bottom: calc(var(--cell-size) * 0.2);
}

.xp-led {
    flex: 1;
    min-width: 0;
    background: #000;
    color: #f00;
    font-family: 'Courier New', Courier, monospace;
    font-size: calc(var(--cell-size) * 0.38);
    font-weight: bold;
    line-height: 1;
    padding: 2px 4px;
    text-align: right;
    letter-spacing: 1px;
}

.xp-face-btn {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    width: calc(var(--cell-size) * 1.1);
    height: calc(var(--cell-size) * 1.1);
    font-size: calc(var(--cell-size) * 0.65);
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

.xp-grid-panel {
    display: grid;
    grid-template-columns: repeat(var(--grid-cols), var(--cell-size));
    width: var(--board-width);
    box-sizing: border-box;
    padding: var(--panel-pad);
    gap: 0;
}

.xp-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: calc(var(--cell-size) * 0.62);
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    box-sizing: border-box;
    cursor: default;
    padding: 0;
}

.xp-cell-empty {
    background: #c0c0c0;
    border-style: solid;
    border-width: max(3px, calc(var(--cell-size) * 0.12));
    border-color: #fff #808080 #808080 #fff;
}

.xp-cell-empty:not(:disabled):active {
    border-width: max(2px, calc(var(--cell-size) * 0.08));
    border-color: #808080;
}

.xp-cell-filled {
    background: #c0c0c0;
    border-style: solid;
    border-width: max(2px, calc(var(--cell-size) * 0.06));
    border-color: #808080;
}

.xp-cell:disabled {
    cursor: default;
}

.xp-cell-x {
    color: #0000ff;
}

.xp-cell-o {
    color: #ff0000;
}

.xp-cell-winning {
    background: #ffff00;
}
</style>
