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

const activePlayer = computed<Player | null>(() => {
    if (isGameOver.value) return null;
    if (isBotThinking.value) return BOT;
    return currentPlayer.value;
});

const resultText = computed(() => {
    if (result.value) return `${result.value.winner} wins!`;
    if (hasDraw.value) return 'Draw game';
    return '';
});

const oPlayerLabel = computed(() => (props.isBot ? 'Computer' : 'Player 2'));

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
        class="ttt-client select-none"
        :class="boardSize === 5 ? 'ttt-size-5' : 'ttt-size-3'"
        :style="{ '--grid-cols': boardSize }"
    >
        <div class="ttt-status">
            <div
                class="ttt-player"
                :class="{ 'ttt-player-active': activePlayer === 'X' }"
            >
                <span class="ttt-mark ttt-mark-x" aria-hidden="true">X</span>
                <span class="ttt-player-name">Player 1</span>
            </div>

            <div class="ttt-center">
                <button
                    type="button"
                    class="ttt-face-btn"
                    title="New game"
                    @click="resetGame"
                    @mousedown="facePressed = true"
                    @mouseup="facePressed = false"
                    @mouseleave="facePressed = false"
                >{{ faceEmoji }}</button>
                <span v-if="resultText" class="ttt-result">{{ resultText }}</span>
                <span v-else-if="isBotThinking" class="ttt-thinking">Thinking…</span>
                <span v-else class="ttt-moves">{{ moveCount }} {{ moveCount === 1 ? 'move' : 'moves' }}</span>
            </div>

            <div
                class="ttt-player"
                :class="{ 'ttt-player-active': activePlayer === 'O' }"
            >
                <span class="ttt-mark ttt-mark-o" aria-hidden="true" />
                <span class="ttt-player-name">{{ oPlayerLabel }}</span>
            </div>
        </div>

        <div class="ttt-board-frame">
            <div class="ttt-board">
                <button
                    v-for="(cell, index) in board"
                    :key="index"
                    type="button"
                    class="ttt-cell"
                    :class="{
                        'ttt-cell-winning': isWinningCell(index),
                        'ttt-cell-x': cell === 'X',
                        'ttt-cell-o': cell === 'O',
                    }"
                    :disabled="
                        cell !== null
                            || isGameOver
                            || isBotThinking
                            || (isBot && currentPlayer !== HUMAN)
                    "
                    @click="handleCellClick(index)"
                >
                    <span v-if="cell === 'X'" class="ttt-mark ttt-mark-x" aria-hidden="true">X</span>
                    <span v-else-if="cell === 'O'" class="ttt-mark ttt-mark-o" aria-hidden="true" />
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.ttt-client {
    --cell-size: 96px;
    --mark-size: calc(var(--cell-size) * 0.55);
    display: flex;
    flex-direction: column;
    width: 100%;
    background: #c0c0c0;
    padding: 12px 14px 6px;
    box-sizing: border-box;
}

.ttt-size-5 {
    --cell-size: 70px;
}

@media (max-width: 640px) {
    .ttt-client {
        padding: 8px 8px 6px;
    }

    .ttt-size-3 {
        --cell-size: min(96px, calc((100vw - 3.5rem) / 3));
    }

    .ttt-size-5 {
        --cell-size: min(68px, calc((100vw - 3.5rem) / 5));
    }

    .ttt-board-frame {
        align-self: center;
        max-width: 100%;
    }
}

.ttt-status {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 10px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
}

.ttt-player {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    min-width: 64px;
    padding: 4px 8px;
    opacity: 0.45;
    transition: opacity 0.15s;
}

.ttt-player-active {
    opacity: 1;
    background: #ece9d8;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
}

.ttt-player-name {
    color: #000;
    white-space: nowrap;
}

.ttt-center {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 3px;
    text-align: center;
    color: #000;
    font-size: 12px;
    min-width: 0;
}

.ttt-face-btn {
    width: 44px;
    height: 44px;
    padding: 0;
    font-size: 28px;
    line-height: 1;
    cursor: default;
    background: #c0c0c0;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
    display: flex;
    align-items: center;
    justify-content: center;
}

.ttt-face-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.ttt-result {
    font-weight: bold;
    color: #000080;
}

.ttt-thinking {
    color: #808080;
    font-style: italic;
}

.ttt-moves {
    color: #404040;
}

.ttt-board-frame {
    align-self: center;
    border: 3px solid;
    border-color: #808080 #fff #fff #808080;
    background: #808080;
    padding: 2px;
}

.ttt-board {
    display: grid;
    grid-template-columns: repeat(var(--grid-cols), var(--cell-size));
    background: #ece9d8;
}

.ttt-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    margin: 0;
    background: transparent;
    border: none;
    border-right: 2px solid #808080;
    border-bottom: 2px solid #808080;
    cursor: default;
    box-sizing: border-box;
}

.ttt-size-3 .ttt-cell:nth-child(3n) {
    border-right: none;
}

.ttt-size-3 .ttt-cell:nth-last-child(-n + 3) {
    border-bottom: none;
}

.ttt-size-5 .ttt-cell:nth-child(5n) {
    border-right: none;
}

.ttt-size-5 .ttt-cell:nth-last-child(-n + 5) {
    border-bottom: none;
}

.ttt-cell:not(:disabled):hover {
    background: rgba(49, 106, 197, 0.08);
}

.ttt-cell:disabled {
    cursor: default;
}

.ttt-cell-winning {
    background: #90ee90;
}

.ttt-mark {
    display: block;
    line-height: 1;
    pointer-events: none;
}

.ttt-mark-x {
    font-size: var(--mark-size);
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    color: #0000ff;
}

.ttt-mark-o {
    width: var(--mark-size);
    height: var(--mark-size);
    border: calc(var(--mark-size) * 0.12) solid #ff0000;
    border-radius: 50%;
    box-sizing: border-box;
}

.ttt-player .ttt-mark-x {
    font-size: 18px;
}

.ttt-player .ttt-mark-o {
    width: 18px;
    height: 18px;
    border-width: 2.5px;
}
</style>
