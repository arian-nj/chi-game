<script setup lang="ts">
import {
  checkWinner,
  createEmptyBoard,
  getBotMove,
  isDraw,
  type BoardSize,
  type Cell,
  type Player,
} from '@/libs/tictactoe';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useTextDirection } from '@/composables/use-text-direction';

const props = defineProps<{
  isBot: boolean;
  boardSize: BoardSize;
}>();

const emit = defineEmits<{
  quit: [];
}>();

const { t } = useI18n();
const { textDir } = useTextDirection();

const HUMAN: Player = 'X';
const BOT: Player = 'O';

const board = ref<Cell[]>(createEmptyBoard(props.boardSize));
const currentPlayer = ref<Player>(HUMAN);
const winningLine = ref<number[]>([]);
const isBotThinking = ref(false);

const result = computed(() => checkWinner(board.value, props.boardSize));
const hasDraw = computed(() => isDraw(board.value, props.boardSize));
const isGameOver = computed(() => result.value !== null || hasDraw.value);

const statusMessage = computed(() => {
  if (result.value) {
    const { winner } = result.value;
    if (props.isBot) {
      return winner === HUMAN ? t('game.youWin') : t('game.botWins');
    }
    return t('game.playerWins', { player: winner });
  }
  if (hasDraw.value) {
    return t('game.draw');
  }
  if (isBotThinking.value) {
    return t('game.botThinking');
  }
  if (props.isBot) {
    return currentPlayer.value === HUMAN ? t('game.yourTurnX') : t('game.botTurnO');
  }
  return t('game.playerTurn', { player: currentPlayer.value });
});

function resetGame() {
  board.value = createEmptyBoard(props.boardSize);
  currentPlayer.value = HUMAN;
  winningLine.value = [];
  isBotThinking.value = false;
}

function isWinningCell(index: number): boolean {
  return winningLine.value.includes(index);
}

function placeMove(index: number, player: Player) {
  if (board.value[index] !== null || isGameOver.value) {
    return;
  }

  board.value[index] = player;

  const gameResult = checkWinner(board.value, props.boardSize);
  if (gameResult) {
    winningLine.value = gameResult.line;
    return;
  }

  if (isDraw(board.value, props.boardSize)) {
    return;
  }

  currentPlayer.value = player === HUMAN ? BOT : HUMAN;
}

function handleCellClick(index: number) {
  if (isGameOver.value || isBotThinking.value) {
    return;
  }
  if (props.isBot && currentPlayer.value !== HUMAN) {
    return;
  }
  if (board.value[index] !== null) {
    return;
  }

  placeMove(index, currentPlayer.value);
}

function runBotMove() {
  if (!props.isBot || isGameOver.value || currentPlayer.value !== BOT) {
    return;
  }

  isBotThinking.value = true;
  window.setTimeout(() => {
    const move = getBotMove([...board.value], props.boardSize, BOT, HUMAN);
    isBotThinking.value = false;

    if (move >= 0) {
      placeMove(move, BOT);
    }
  }, 350);
}

watch(
  () => [currentPlayer.value, isGameOver.value, props.isBot] as const,
  () => runBotMove(),
  { flush: 'post' },
);

watch(
  () => props.boardSize,
  () => resetGame(),
);
</script>

<template>
  <div class="flex flex-col items-center gap-5">
    <div
      :dir="textDir"
      class="w-full rounded-2xl border border-white/10 bg-custom-lite-blue/40 p-4 shadow-md"
      role="status"
      aria-live="polite"
    >
      <p class="text-center text-lg font-bold text-white">{{ statusMessage }}</p>
    </div>

    <div
      class="grid w-full max-w-md gap-2"
      :class="boardSize === 3 ? 'grid-cols-3' : 'grid-cols-5'"
      role="grid"
      :aria-label="t('game.tttBoardAria', { size: boardSize })"
    >
      <button
        v-for="(cell, index) in board"
        :key="index"
        type="button"
        class="aspect-square rounded-xl border border-white/15 bg-custom-deep-blue/80 font-extrabold transition duration-200 select-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 disabled:cursor-not-allowed"
        :class="[
          boardSize === 3 ? 'text-4xl sm:text-5xl' : 'text-2xl sm:text-3xl',
          cell === 'X' ? 'text-sky-300' : cell === 'O' ? 'text-rose-300' : 'text-transparent',
          isWinningCell(index) ? 'bg-green-500/30 ring-2 ring-green-400' : 'hover:bg-white/10',
          isGameOver || isBotThinking ? 'cursor-default' : 'cursor-pointer',
        ]"
        :disabled="cell !== null || isGameOver || isBotThinking || (isBot && currentPlayer !== HUMAN)"
        :aria-label="
          cell
            ? t('game.cellOccupied', { index: index + 1, symbol: cell })
            : t('game.cellEmpty', { index: index + 1 })
        "
        @click="handleCellClick(index)"
      >
        {{ cell ?? '' }}
      </button>
    </div>

    <div class="flex w-full max-w-md flex-col gap-3 sm:flex-row">
      <button
        type="button"
        class="flex-1 rounded-xl bg-white/90 px-4 py-3 text-lg font-bold text-custom-blue transition hover:bg-white"
        @click="resetGame"
      >
        {{ t('game.newGame') }}
      </button>
      <button
        type="button"
        class="flex-1 rounded-xl border border-white/20 bg-transparent px-4 py-3 text-lg font-bold text-blue-100 transition hover:bg-white/10"
        @click="emit('quit')"
      >
        {{ t('game.backToSettings') }}
      </button>
    </div>
  </div>
</template>
