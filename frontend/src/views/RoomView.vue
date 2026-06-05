<script setup lang="ts">
import { useTextDirection } from '@/composables/use-text-direction';
import { InviteService } from '@/gen/invite/v1/invite_pb';
import { createApiClient } from '@/libs/api-client';
import { gamesData } from '@/libs/game';
import { leaveCurrentRoom } from '@/libs/invite-room';
import { roomLobbyPath } from '@/libs/room-url';
import { useQuery } from '@tanstack/vue-query';
import { computed, watch } from 'vue';
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { textDir } = useTextDirection();

const locale = computed(() => route.params.locale as string);
const inviteCode = computed(() => {
  const raw = route.params.code;
  if (typeof raw !== 'string') {
    return '';
  }
  return raw.trim().toUpperCase();
});

const enabledGames = computed(() => gamesData.filter(game => game.isEnable));

const client = createApiClient(InviteService);

const { data, isError } = useQuery({
  queryKey: computed(() => ['invite-room', inviteCode.value]),
  queryFn: ({ signal }) => client.getInviteRoom({ inviteCode: inviteCode.value }, { signal }),
  refetchInterval: 2000,
  enabled: computed(() => Boolean(inviteCode.value)),
});

const players = computed(() => data.value?.players ?? []);

function playerLabel(displayName: string, username: string) {
  const name = displayName?.trim();
  if (name) {
    return name;
  }
  return username ? `@${username}` : t('invite.unknownPlayer');
}

watch(
  players,
  (list) => {
    if (!inviteCode.value) {
      return;
    }
    if (list.length < 2 && route.name === 'room-play') {
      void router.replace({
        name: 'room-code',
        params: { locale: locale.value, code: inviteCode.value },
      });
    }
  },
  { flush: 'post' },
);

onBeforeRouteLeave((to) => {
  if (to.name === 'room' || to.name === 'room-code' || to.name === 'room-play') {
    return;
  }
  void leaveCurrentRoom(inviteCode.value);
});
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center pb-3 pt-14 text-white relative">
    <!-- Home Button -->
    <RouterLink
      :to="`/${locale}`"
      class="absolute left-4 top-14 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :aria-label="t('nav.backToHome')"
    >
      <span aria-hidden="true">←</span>
      {{ t('nav.home') }}
    </RouterLink>

    <!-- Play Title -->
    <h1 class="text-4xl font-bold text-white mb-4 mt-2 animate-pop select-none drop-shadow-sm uppercase">
      {{ t('invite.playTitle') }}
    </h1>
    <!-- Players -->
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

      <!-- Choose Game -->
      <div class="flex flex-col gap-3" :dir="textDir">
        <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.chooseGame') }}</span>
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <button
            v-for="game in enabledGames"
            :key="game.key"
            @click="(game.key)"
            class="rounded-xl border border-white/20 bg-green-500/90 px-4 py-4 text-center text-lg font-bold text-white shadow-md transition hover:bg-green-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
          >
            {{ t(`games.${game.key}`) }}
          </button>
        </div>
      </div>

      <!-- Back to Lobby -->
      <RouterLink
        :to="roomLobbyPath(locale, inviteCode)"
        class="text-center text-sm font-semibold text-blue-100 underline underline-offset-4 hover:text-white"
      >
        {{ t('invite.backToLobby') }}
      </RouterLink>
    </div>
  </div>
</template>
