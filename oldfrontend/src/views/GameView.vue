<script setup lang="ts">
import Connect4Game from '@/components/connect4/Connect4Game.vue';
import TicTacToeGame from '@/components/tictactoe/TicTacToeGame.vue';
import { useToast } from '@/components/Toast.vue';
import type { BotDifficulty } from '@/libs/bot-difficulty';
import type { Connect4Settings } from '@/libs/connect4';
import { gamesData, type GameKey } from '@/libs/game';
import type { TicTacToeSettings } from '@/libs/tictactoe';
import GamePlayView, { type PlaySettings } from '@/views/GamePlayView.vue';
import GameRulesView from '@/views/GameRulesView.vue';
import { computed, ref, useTemplateRef } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

type ActiveGame =
  | { type: 'tic-tac-toe'; settings: TicTacToeSettings }
  | { type: 'connect-4'; settings: Connect4Settings };

const route = useRoute();
const { t } = useI18n();
const locale = computed(() => route.params.locale as string);
const urlGameName = route.params.game as GameKey;
const gameData = gamesData.find(game => game.key === urlGameName);
const isTicTacToe = computed(() => route.params.game === 'tic-tac-toe');
const isConnect4 = computed(() => route.params.game === 'connect-4');
const isEnabled = computed(() => Boolean(gameData?.isEnable));
const isSupported = computed(() => isTicTacToe.value || isConnect4.value);
const isPlayable = computed(() => isEnabled.value && isSupported.value);
const isPlaying = ref(false);
const activeGame = ref<ActiveGame | null>(null);
const playButtonAnimationKey = ref(0);
const playViewRef = useTemplateRef<{ getSettings: () => PlaySettings }>('playViewRef');

function playGame() {
  const toast = useToast();

  if (!isPlayable.value) {
    toast.toast.info(t('play.comingSoon'));
    return;
  }

  const playView = playViewRef.value;
  const settings =
    playView && typeof playView.getSettings === 'function'
      ? playView.getSettings()
      : { isBot: true, botDifficulty: 'medium' as BotDifficulty };

  if (isTicTacToe.value) {
    activeGame.value = {
      type: 'tic-tac-toe',
      settings: {
        isBot: settings.isBot,
        boardSize: 'boardSize' in settings ? settings.boardSize : 3,
        botDifficulty: settings.botDifficulty,
      },
    };
  } else if (isConnect4.value) {
    activeGame.value = {
      type: 'connect-4',
      settings: { isBot: settings.isBot, botDifficulty: settings.botDifficulty },
    };
  }

  isPlaying.value = true;
}

function quitGame() {
  isPlaying.value = false;
  activeGame.value = null;
  playButtonAnimationKey.value += 1;
}
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center pb-3 pt-14 text-white relative">
    <RouterLink
      :to="`/${locale}`"
      class="absolute left-4 top-14 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :aria-label="t('nav.backToHome')"
    >
      <span aria-hidden="true">←</span>
      {{ t('nav.home') }}
    </RouterLink>

    <h1 class="text-4xl font-bold text-white mb-8 mt-2 animate-pop select-none drop-shadow-sm uppercase">
      {{ gameData ? t(`games.${gameData.key}`) : '' }}
    </h1>

    <div class="w-full max-w-3xl px-4 flex-1 pb-28">
      <TicTacToeGame
        v-if="isPlaying && activeGame?.type === 'tic-tac-toe'"
        :is-bot="activeGame.settings.isBot"
        :board-size="activeGame.settings.boardSize"
        :bot-difficulty="activeGame.settings.botDifficulty"
        @quit="quitGame"
      />
      <Connect4Game
        v-else-if="isPlaying && activeGame?.type === 'connect-4'"
        :is-bot="activeGame.settings.isBot"
        :bot-difficulty="activeGame.settings.botDifficulty"
        @quit="quitGame"
      />
      <div v-else class="flex flex-col gap-6">
        <GamePlayView ref="playViewRef" />
        <GameRulesView />
      </div>
    </div>

    <div
      v-show="!isPlaying && isPlayable"
      :key="playButtonAnimationKey"
      class="animate-play-btn-enter fixed bottom-10 left-1/2 z-20 w-3/4 max-w-md"
    >
      <button
        type="button"
        class="w-full rounded-lg bg-green-500 p-3 text-3xl font-extrabold text-white"
        @click="playGame"
      >
        <span class="inline-block animate-bounce text-3xl">🚀</span><span>{{ t('play.button') }}</span>
      </button>
    </div>
  </div>
</template>
