<script setup lang="ts">
import { useRoomSession } from '@/composables/use-room-session';
import { gamesData } from '@/libs/game';
import { roomLobbyPath } from '@/libs/room-url';
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const router = useRouter();

const {
  locale,
  roomCode,
  players,
  isError,
  textDir,
  playerLabel,
} = useRoomSession();

const enabledGames = computed(() => gamesData.filter(game => game.isEnable));

function chooseGame(gameKey: string) {
  void router.push({ name: 'game', params: { locale: locale.value, game: gameKey } });
}
</script>

<template>
  <h1 class="text-4xl font-bold text-white mb-4 mt-2 animate-pop select-none drop-shadow-sm uppercase">
    {{ t('invite.playTitle') }}
  </h1>

  <div class="w-full max-w-3xl px-4 flex-1 pb-10 flex flex-col gap-6">
    <div
      class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md flex flex-col gap-4 p-5"
      :dir="textDir"
    >
      <p class="text-sm text-blue-100/80">{{ t('invite.playHint') }}</p>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.playersLabel') }}</span>
        <p v-if="isError" class="text-sm text-red-300">{{ t('invite.loadError') }}</p>
        <ul v-else class="flex flex-col gap-2">
          <li
            v-for="player in players"
            :key="String(player.id)"
            class="rounded-xl border border-white/10 bg-custom-deep-blue/60 px-4 py-3 text-lg font-semibold text-white"
          >
            {{ playerLabel(player.displayName, player.username) }}
          </li>
        </ul>
      </div>
    </div>

    <div class="flex flex-col gap-3" :dir="textDir">
      <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.chooseGame') }}</span>
      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <button
          v-for="game in enabledGames"
          :key="game.key"
          type="button"
          class="rounded-xl border border-white/20 bg-green-500/90 px-4 py-4 text-center text-lg font-bold text-white shadow-md transition hover:bg-green-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
          @click="chooseGame(game.key)"
        >
          {{ t(`games.${game.key}`) }}
        </button>
      </div>
    </div>

    <RouterLink
      :to="roomLobbyPath(locale, roomCode)"
      class="text-center text-sm font-semibold text-blue-100 underline underline-offset-4 hover:text-white"
    >
      {{ t('invite.backToLobby') }}
    </RouterLink>
  </div>
</template>
