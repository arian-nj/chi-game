<script setup lang="ts">
import TicTacToeGame from '@/components/tictactoe/TicTacToeGame.vue';
import { useToast } from '@/components/Toast.vue';
import type { TicTacToeSettings } from '@/libs/tictactoe';
import { computed, ref, useTemplateRef } from 'vue';
import { useRoute } from 'vue-router';
import { gamesData } from '../libs/game';

const route = useRoute();
const urlGameName = route.params.game;
const gameData = gamesData.find(game => game.key === urlGameName);
const isTicTacToe = computed(() => route.params.game === 'tictactoe');
const isPlayTab = computed(() => route.name === 'game-play' || route.path.match(/\/game\/[^/]+$/));
const isPlaying = ref(false);
const gameSettings = ref<TicTacToeSettings | null>(null);
const playButtonAnimationKey = ref(0);
const playViewRef = useTemplateRef<{ getSettings: () => TicTacToeSettings }>('playViewRef');

// Tabs
const tabs = [
    { name: 'Play', icon: '🚀', to: `/game/${urlGameName}` },
    { name: 'Rules', icon: '📖', to: `/game/${urlGameName}/rules` },
];

function playGame() {
  const toast = useToast();

  if (!isTicTacToe.value) {
    toast.toast.info('This game is coming soon');
    return;
  }

  const playView = playViewRef.value;
  gameSettings.value =
    playView && typeof playView.getSettings === 'function'
      ? playView.getSettings()
      : { isBot: true, boardSize: 3 };
  isPlaying.value = true;
}

function quitGame() {
  isPlaying.value = false;
  gameSettings.value = null;
  playButtonAnimationKey.value += 1;
}

</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center pb-3 text-white relative">
    <RouterLink
      to="/"
      class="absolute left-4 top-4 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      aria-label="Back to home"
    >
      <span aria-hidden="true">←</span>
      Home
    </RouterLink>

    <h1 class="text-4xl font-bold text-white mb-8 mt-2 animate-pop select-none drop-shadow-sm uppercase">
      {{ gameData?.name }}
    </h1>
    
    <!-- Tab Navigation -->
    <div class="flex space-x-4 bg-custom-lite-blue/70 rounded-xl mb-6 px-2 py-1 shadow-lg">
        <RouterLink
            v-for="tab in tabs"
            :key="tab.name"
            :to="tab.to"
            custom
            v-slot="{ isExactActive, navigate }"
        >

            <a
            :href="tab.to"
            @click="navigate"
            :class="[
                'px-6 py-2 rounded-lg transition font-bold flex items-center gap-2 cursor-pointer',
                isExactActive /* USING IT HERE NOW */
                ? 'bg-white/90 text-custom-blue shadow'
                : 'bg-transparent text-blue-100 hover:bg-white/30'
            ]"
            >
            <span>{{ tab.icon }}</span> {{ tab.name }}
            </a>
        </RouterLink>
    </div>

    <div class="w-full max-w-3xl px-4 flex-1 pb-28">
      <TicTacToeGame
        v-if="isPlaying && gameSettings"
        :is-bot="gameSettings.isBot"
        :board-size="gameSettings.boardSize"
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
        <span class="inline-block animate-bounce text-3xl">🚀</span><span>Play</span>
      </button>
    </div>
  </div>
</template>