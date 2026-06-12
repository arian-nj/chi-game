<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import {
    CONNECT4_COLS,
    CONNECT4_ROWS,
    checkWinnerFrom,
    cloneBoard,
    createEmptyBoard,
    dropDisc,
    getBotMove,
    getDropRow,
    isBoardFull,
    isWinningCell,
    type Board,
    type Connect4Result,
    type Player,
} from '@/lib/connect4/connect4';
import type { BotDifficulty } from '@/lib/connect4/types';

const props = defineProps<{
    isBot: boolean;
    botDifficulty: BotDifficulty;
}>();

const HUMAN: Player = 'R';
const BOT: Player = 'Y';

const board = ref<Board>(createEmptyBoard());
const currentPlayer = ref<Player>(HUMAN);
const winningCells = ref<Connect4Result['cells']>([]);
const isBotThinking = ref(false);
const facePressed = ref(false);
const moveCount = ref(0);

const result = computed(() => {
    if (winningCells.value.length > 0) {
        const winner = board.value[winningCells.value[0]!.row]![winningCells.value[0]!.col]!;
        if (winner === 'R' || winner === 'Y') {
            return { winner, cells: winningCells.value };
        }
    }
    return null;
});

const hasDraw = computed(() => !result.value && isBoardFull(board.value));
const isGameOver = computed(() => result.value !== null || hasDraw.value);

const activePlayer = computed<Player | null>(() => {
    if (isGameOver.value) return null;
    if (isBotThinking.value) return BOT;
    return currentPlayer.value;
});

const resultText = computed(() => {
    if (result.value) {
        if (props.isBot) {
            return result.value.winner === HUMAN ? 'You win!' : 'Computer wins!';
        }
        return result.value.winner === HUMAN ? 'Red wins!' : 'Yellow wins!';
    }
    if (hasDraw.value) return 'Draw game';
    return '';
});

const yellowPlayerLabel = computed(() => (props.isBot ? 'Computer' : 'Player 2'));

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
    board.value = createEmptyBoard();
    currentPlayer.value = HUMAN;
    winningCells.value = [];
    isBotThinking.value = false;
    facePressed.value = false;
    moveCount.value = 0;
}

function canPlayColumn(col: number): boolean {
    if (isGameOver.value || isBotThinking.value) return false;
    if (props.isBot && currentPlayer.value !== HUMAN) return false;
    return getDropRow(board.value, col) !== -1;
}

function placeInColumn(col: number, player: Player) {
    const row = dropDisc(board.value, col, player);
    if (row === -1) return;

    moveCount.value++;

    const gameResult = checkWinnerFrom(board.value, row, col);
    if (gameResult) {
        winningCells.value = gameResult.cells;
        return;
    }

    if (isBoardFull(board.value)) return;

    currentPlayer.value = player === HUMAN ? BOT : HUMAN;
}

function handleColumnClick(col: number) {
    if (!canPlayColumn(col)) return;
    placeInColumn(col, currentPlayer.value);
}

function runBotMove() {
    if (!props.isBot || isGameOver.value || currentPlayer.value !== BOT) return;

    isBotThinking.value = true;
    window.setTimeout(() => {
        const col = getBotMove(cloneBoard(board.value), BOT, HUMAN, props.botDifficulty);
        isBotThinking.value = false;
        if (col >= 0) placeInColumn(col, BOT);
    }, 400);
}

watch(
    () => [currentPlayer.value, isGameOver.value, props.isBot] as const,
    () => runBotMove(),
    { flush: 'post' },
);

watch(
    () => [props.isBot, props.botDifficulty],
    () => resetGame(),
);

defineExpose({ resetGame });
</script>

<template>
    <div class="c4-client select-none">
        <div class="c4-status">
            <div
                class="c4-player"
                :class="{ 'c4-player-active': activePlayer === 'R' }"
            >
                <span class="c4-disc c4-disc-red" aria-hidden="true" />
                <span class="c4-player-name">Player 1</span>
            </div>

            <div class="c4-center">
                <button
                    type="button"
                    class="c4-face-btn"
                    title="New game"
                    @click="resetGame"
                    @mousedown="facePressed = true"
                    @mouseup="facePressed = false"
                    @mouseleave="facePressed = false"
                >{{ faceEmoji }}</button>
                <span v-if="resultText" class="c4-result">{{ resultText }}</span>
                <span v-else-if="isBotThinking" class="c4-thinking">Thinking…</span>
                <span v-else class="c4-moves">{{ moveCount }} {{ moveCount === 1 ? 'move' : 'moves' }}</span>
            </div>

            <div
                class="c4-player"
                :class="{ 'c4-player-active': activePlayer === 'Y' }"
            >
                <span class="c4-disc c4-disc-yellow" aria-hidden="true" />
                <span class="c4-player-name">{{ yellowPlayerLabel }}</span>
            </div>
        </div>

        <div class="c4-board-frame">
            <div class="c4-drop-row">
                <button
                    v-for="col in CONNECT4_COLS"
                    :key="`drop-${col}`"
                    type="button"
                    class="c4-drop-btn"
                    :disabled="!canPlayColumn(col - 1)"
                    :aria-label="`Drop disc in column ${col}`"
                    @click="handleColumnClick(col - 1)"
                >↓</button>
            </div>

            <div class="c4-board" role="grid" aria-label="Connect 4 board">
                <template v-for="row in CONNECT4_ROWS" :key="`row-${row}`">
                    <button
                        v-for="col in CONNECT4_COLS"
                        :key="`cell-${row}-${col}`"
                        type="button"
                        class="c4-cell"
                        :class="{
                            'c4-cell-winning': isWinningCell(winningCells, row - 1, col - 1),
                        }"
                        :disabled="!canPlayColumn(col - 1)"
                        :aria-label="`Column ${col}, row ${row}`"
                        @click="handleColumnClick(col - 1)"
                    >
                        <span
                            class="c4-disc"
                            :class="{
                                'c4-disc-red': board[row - 1]?.[col - 1] === 'R',
                                'c4-disc-yellow': board[row - 1]?.[col - 1] === 'Y',
                                'c4-disc-empty': board[row - 1]?.[col - 1] === null,
                            }"
                        />
                    </button>
                </template>
            </div>
        </div>
    </div>
</template>

<style scoped>
.c4-client {
    --cell-size: 58px;
    --disc-size: 48px;
    width: fit-content;
    background: #c0c0c0;
    padding: 12px 14px 6px;
    box-sizing: border-box;
}

@media (max-width: 640px) {
    .c4-client {
        width: 100%;
        --cell-size: min(52px, calc((100vw - 3.5rem) / 7));
        --disc-size: calc(var(--cell-size) * 0.82);
        padding: 8px 8px 6px;
    }

    .c4-board-frame {
        width: fit-content;
        max-width: 100%;
        margin: 0 auto;
    }
}

.c4-status {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 10px;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    font-size: 11px;
}

.c4-player {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
    min-width: 64px;
    padding: 4px 8px;
    opacity: 0.45;
    transition: opacity 0.15s;
}

.c4-player-active {
    opacity: 1;
    background: #ece9d8;
    border: 2px solid;
    border-color: #fff #808080 #808080 #fff;
}

.c4-player-name {
    color: #000;
    white-space: nowrap;
}

.c4-center {
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

.c4-face-btn {
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

.c4-face-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.c4-result {
    font-weight: bold;
    color: #000080;
}

.c4-thinking {
    color: #808080;
    font-style: italic;
}

.c4-moves {
    color: #404040;
}

.c4-board-frame {
    border: 3px solid;
    border-color: #808080 #fff #fff #808080;
    background: #000080;
    padding: 4px;
}

.c4-drop-row {
    display: grid;
    grid-template-columns: repeat(7, var(--cell-size));
    margin-bottom: 2px;
}

.c4-drop-btn {
    width: var(--cell-size);
    height: 20px;
    padding: 0;
    font-size: 12px;
    line-height: 1;
    color: #fff;
    background: transparent;
    border: none;
    cursor: default;
    opacity: 0.7;
}

.c4-drop-btn:not(:disabled):hover {
    opacity: 1;
    background: rgba(255, 255, 255, 0.15);
}

.c4-drop-btn:disabled {
    opacity: 0.25;
    cursor: default;
}

.c4-board {
    display: grid;
    grid-template-columns: repeat(7, var(--cell-size));
    gap: 2px;
}

.c4-cell {
    width: var(--cell-size);
    height: var(--cell-size);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    margin: 0;
    background: #000080;
    border: none;
    cursor: default;
    box-sizing: border-box;
}

.c4-cell:not(:disabled):hover {
    background: #0000a0;
}

.c4-cell-winning {
    background: #006400;
}

.c4-disc {
    display: block;
    width: var(--disc-size);
    height: var(--disc-size);
    border-radius: 50%;
    box-sizing: border-box;
    pointer-events: none;
}

.c4-disc-red {
    background: radial-gradient(circle at 35% 30%, #ff6666, #cc0000 60%, #990000);
    border: 1px solid #880000;
}

.c4-disc-yellow {
    background: radial-gradient(circle at 35% 30%, #ffff99, #ffcc00 60%, #cc9900);
    border: 1px solid #997700;
}

.c4-disc-empty {
    background: #ece9d8;
    border: 1px solid #808080;
}

.c4-player .c4-disc {
    width: 18px;
    height: 18px;
}
</style>
