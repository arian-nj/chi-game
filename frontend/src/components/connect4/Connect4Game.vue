<script setup lang="ts">
import type { BotDifficulty } from '@/libs/bot-difficulty';
import {
  CONNECT4_COLS,
  CONNECT4_ROWS,
  checkWinnerFrom,
  createEmptyBoard,
  dropDisc,
  getBotMove,
  getDropRow,
  isBoardFull,
  isWinningCell,
  type Board,
  type Connect4Result,
  type Player,
} from '@/libs/connect4';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useTextDirection } from '@/composables/use-text-direction';
import GameTurnIndicator, { type TurnPlayer } from '@/components/GameTurnIndicator.vue';

const props = defineProps<{
  isBot: boolean;
  botDifficulty: BotDifficulty;
}>();

const emit = defineEmits<{
  quit: [];
}>();

const { t } = useI18n();
const { textDir } = useTextDirection();

const HUMAN: Player = 'R';
const BOT: Player = 'Y';

const board = ref<Board>(createEmptyBoard());
const currentPlayer = ref<Player>(HUMAN);
const winningCells = ref<Connect4Result['cells']>([]);
const isBotThinking = ref(false);
const lastDropCol = ref<number | null>(null);

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

const statusMessage = computed(() => {
  if (result.value) {
    const { winner } = result.value;
    if (props.isBot) {
      return winner === HUMAN ? t('game.youWin') : t('game.botWins');
    }
    return winner === HUMAN ? t('game.redWins') : t('game.yellowWins');
  }
  if (hasDraw.value) {
    return t('game.draw');
  }
  if (isBotThinking.value) {
    return t('game.botThinking');
  }
  if (props.isBot) {
    return currentPlayer.value === HUMAN ? t('game.yourTurnRed') : t('game.botTurnYellow');
  }
  return currentPlayer.value === HUMAN ? t('game.redTurn') : t('game.yellowTurn');
});

const turnPlayers = computed<TurnPlayer[]>(() => [
  {
    key: HUMAN,
    label: props.isBot ? t('game.you') : t('game.red'),
    markerType: 'disc',
    markerClass: 'bg-rose-400',
  },
  {
    key: BOT,
    label: props.isBot ? t('game.botPlayer') : t('game.yellow'),
    markerType: 'disc',
    markerClass: 'bg-yellow-300',
  },
]);

const activePlayerKey = computed(() => {
  if (result.value) {
    return result.value.winner;
  }
  if (hasDraw.value) {
    return null;
  }
  if (isBotThinking.value) {
    return BOT;
  }
  return currentPlayer.value;
});

const turnStatus = computed(() => {
  if (result.value) {
    return 'win' as const;
  }
  if (hasDraw.value) {
    return 'draw' as const;
  }
  if (isBotThinking.value) {
    return 'thinking' as const;
  }
  return 'playing' as const;
});

function resetGame() {
  board.value = createEmptyBoard();
  currentPlayer.value = HUMAN;
  winningCells.value = [];
  isBotThinking.value = false;
  lastDropCol.value = null;
}

function canPlayColumn(col: number): boolean {
  if (isGameOver.value || isBotThinking.value) {
    return false;
  }
  if (props.isBot && currentPlayer.value !== HUMAN) {
    return false;
  }
  return getDropRow(board.value, col) !== -1;
}

function placeInColumn(col: number, player: Player) {
  const row = dropDisc(board.value, col, player);
  if (row === -1) {
    return;
  }

  lastDropCol.value = col;

  const gameResult = checkWinnerFrom(board.value, row, col);
  if (gameResult) {
    winningCells.value = gameResult.cells;
    return;
  }

  if (isBoardFull(board.value)) {
    return;
  }

  currentPlayer.value = player === HUMAN ? BOT : HUMAN;
}

function handleColumnClick(col: number) {
  if (!canPlayColumn(col)) {
    return;
  }
  placeInColumn(col, currentPlayer.value);
}

function runBotMove() {
  if (!props.isBot || isGameOver.value || currentPlayer.value !== BOT) {
    return;
  }

  isBotThinking.value = true;
  window.setTimeout(() => {
    const col = getBotMove(cloneBoardForBot(), BOT, HUMAN, props.botDifficulty);
    isBotThinking.value = false;

    if (col >= 0) {
      placeInColumn(col, BOT);
    }
  }, 400);
}

function cloneBoardForBot(): Board {
  return board.value.map(row => [...row]);
}

watch(
  () => [currentPlayer.value, isGameOver.value, props.isBot] as const,
  () => runBotMove(),
  { flush: 'post' },
);
</script>

<template>
  <div class="flex flex-col items-center gap-5">
    <GameTurnIndicator
      :dir="textDir"
      :message="statusMessage"
      :status="turnStatus"
      :players="turnPlayers"
      :active-player-key="activePlayerKey"
      :show-players="!hasDraw"
    />

    <div class="w-full max-w-lg rounded-2xl border border-white/10 bg-custom-deep-blue/90 p-3 shadow-lg">
      <div class="mb-2 grid grid-cols-7 gap-1.5">
        <button
          v-for="col in CONNECT4_COLS"
          :key="`drop-${col}`"
          type="button"
          class="rounded-lg py-1 text-lg font-bold text-blue-100/70 transition hover:bg-white/10 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 disabled:cursor-not-allowed disabled:opacity-40"
          :disabled="!canPlayColumn(col - 1)"
          :aria-label="t('game.dropDiscColumn', { col })"
          @click="handleColumnClick(col - 1)"
        >
          ↓
        </button>
      </div>

      <div
        class="grid grid-cols-7 gap-1.5 rounded-xl bg-custom-blue/80 p-2"
        role="grid"
        :aria-label="t('game.connect4BoardAria')"
      >
        <template v-for="row in CONNECT4_ROWS" :key="`row-${row}`">
          <button
            v-for="col in CONNECT4_COLS"
            :key="`cell-${row}-${col}`"
            type="button"
            class="aspect-square rounded-full border border-white/10 bg-custom-lite-blue/50 p-1 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 disabled:cursor-default"
            :class="[
              canPlayColumn(col - 1) ? 'cursor-pointer hover:bg-white/10' : 'cursor-default',
              lastDropCol === col - 1 && board[row - 1]?.[col - 1] ? 'ring-2 ring-white/30' : '',
            ]"
            :disabled="!canPlayColumn(col - 1)"
            :aria-label="t('game.columnRow', { col, row })"
            @click="handleColumnClick(col - 1)"
          >
            <span
              class="block size-full rounded-full transition-transform duration-300"
              :class="[
                board[row - 1]?.[col - 1] === 'R'
                  ? 'scale-100 bg-rose-400 shadow-inner'
                  : board[row - 1]?.[col - 1] === 'Y'
                    ? 'scale-100 bg-yellow-300 shadow-inner'
                    : 'scale-0 bg-transparent',
                isWinningCell(winningCells, row - 1, col - 1)
                  ? 'ring-4 ring-green-400 shadow-[0_0_14px_3px_rgba(74,222,128,0.55)]'
                  : '',
              ]"
            ></span>
          </button>
        </template>
      </div>
    </div>

    <div class="flex w-full max-w-lg flex-col gap-3 sm:flex-row">
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
