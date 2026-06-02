<script setup lang="ts">
import Connect4Game from '@/components/connect4/Connect4Game.vue';
import TicTacToeGame from '@/components/tictactoe/TicTacToeGame.vue';
import { useToast } from '@/components/Toast.vue';
import type { BotDifficulty } from '@/libs/bot-difficulty';
import type { Connect4Settings } from '@/libs/connect4';
import type { TicTacToeSettings } from '@/libs/tictactoe';
import { gamesData, type GameKey } from '@/libs/game';
import type { PlaySettings } from '@/views/GamePlayView.vue';
import { computed, ref, useTemplateRef } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';

type ActiveGame =
  | { type: 'tictactoe'; settings: TicTacToeSettings }
  | { type: 'conn4'; settings: Connect4Settings };

const route = useRoute();
const { t } = useI18n();
const urlGameName = route.params.game as GameKey;
const gameData = gamesData.find(game => game.key === urlGameName);
const isTicTacToe = computed(() => route.params.game === 'tictactoe');
const isConnect4 = computed(() => route.params.game === 'conn4');
const isPlayable = computed(() => isTicTacToe.value || isConnect4.value);
const isPlayTab = computed(() => route.name === 'game-play' || route.path.match(/\/game\/[^/]+$/));
const isPlaying = ref(false);
const activeGame = ref<ActiveGame | null>(null);
const playButtonAnimationKey = ref(0);
const playViewRef = useTemplateRef<{ getSettings: () => PlaySettings }>('playViewRef');

const tabs = computed(() => [
  { name: t('nav.play'), icon: '🚀', to: `/game/${urlGameName}` },
  { name: t('nav.rules'), icon: '📖', to: `/game/${urlGameName}/rules` },
]);

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
      type: 'tictactoe',
      settings: {
        isBot: settings.isBot,
        boardSize: 'boardSize' in settings ? settings.boardSize : 3,
        botDifficulty: settings.botDifficulty,
      },
    };
  } else if (isConnect4.value) {
    activeGame.value = {
      type: 'conn4',
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
      to="/"
      class="absolute left-4 top-14 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :aria-label="t('nav.backToHome')"
    >
      <span aria-hidden="true">←</span>
      {{ t('nav.home') }}
    </RouterLink>

    <h1 class="text-4xl font-bold text-white mb-8 mt-2 animate-pop select-none drop-shadow-sm uppercase">
      {{ gameData ? t(`games.${gameData.key}`) : '' }}
    </h1>

    <div
      v-show="!isPlaying"
      class="flex gap-4 bg-custom-lite-blue/70 rounded-xl mb-6 px-2 py-1 shadow-lg"
    >
      <RouterLink
        v-for="tab in tabs"
        :key="tab.to"
        :to="tab.to"
        custom
        v-slot="{ isExactActive, navigate }"
      >
        <a
          :href="tab.to"
          @click="navigate"
          :class="[
            'px-6 py-2 rounded-lg transition font-bold flex items-center gap-2 cursor-pointer',
            isExactActive ? 'bg-white/90 text-custom-blue shadow' : 'bg-transparent text-blue-100 hover:bg-white/30',
          ]"
        >
          <span>{{ tab.icon }}</span> {{ tab.name }}
        </a>
      </RouterLink>
    </div>

    <div class="w-full max-w-3xl px-4 flex-1 pb-28">
      <TicTacToeGame
        v-if="isPlaying && activeGame?.type === 'tictactoe'"
        :is-bot="activeGame.settings.isBot"
        :board-size="activeGame.settings.boardSize"
        :bot-difficulty="activeGame.settings.botDifficulty"
        @quit="quitGame"
      />
      <Connect4Game
        v-else-if="isPlaying && activeGame?.type === 'conn4'"
        :is-bot="activeGame.settings.isBot"
        :bot-difficulty="activeGame.settings.botDifficulty"
        @quit="quitGame"
      />
      <RouterView v-else v-slot="{ Component }">
        <component :is="Component" ref="playViewRef" />
      </RouterView>
    </div>

    <div
      v-show="isPlayTab && !isPlaying"
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
