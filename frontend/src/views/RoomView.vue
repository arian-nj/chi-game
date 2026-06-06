<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { useTextDirection } from '@/composables/use-text-direction';
import { RoomService } from '@/gen/room/v1/room_pb';
import { createApiClient } from '@/libs/api-client';
import { gamesData } from '@/libs/game';
import { joinRoomWithCode, leaveCurrentRoom } from '@/libs/room-api';
import { roomLobbyPath } from '@/libs/room-url';
import { ConnectError } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { computed, ref, watch } from 'vue';
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { textDir } = useTextDirection();
const { isGuest } = useGuestAuth();

const locale = computed(() => route.params.locale as string);
const roomCode = computed(() => {
  const raw = route.params.code;
  if (typeof raw !== 'string') {
    return '';
  }
  return raw.trim().toUpperCase();
});

const joinAttempted = ref(false);
const joinReady = ref(false);

const enabledGames = computed(() => gamesData.filter(game => game.isEnable));

const client = createApiClient(RoomService);

const { data, isError } = useQuery({
  queryKey: computed(() => ['room', roomCode.value]),
  queryFn: ({ signal }) => client.getRoom({ code: roomCode.value }, { signal }),
  refetchInterval: 2000,
  enabled: computed(() => Boolean(roomCode.value) && joinReady.value),
});

const players = computed(() => data.value?.players ?? []);

function inviteErrorMessage(err: unknown): string {
  if (err instanceof ConnectError) {
    const msg = err.message.toLowerCase();
    if (msg.includes('full')) {
      return t('invite.roomFull');
    }
    if (msg.includes('invalid') || msg.includes('expired') || msg.includes('not found')) {
      return t('invite.invalidCode');
    }
    if (err.code === 16) {
      return t('invite.needAuth');
    }
  }
  return '';
}

function playerLabel(displayName: string, username: string) {
  const name = displayName?.trim();
  if (name) {
    return name;
  }
  return username ? `@${username}` : t('invite.unknownPlayer');
}

async function ensureJoined(code: string) {
  try {
    await joinRoomWithCode(code);
    joinReady.value = true;
  } catch (err) {
    const toast = useToast();
    const specific = inviteErrorMessage(err);
    toast.toast.error(specific || t('invite.joinFailed'));
    await router.replace({
      name: 'room-code',
      params: { locale: locale.value, code },
    });
  }
}

watch(
  roomCode,
  async (code) => {
    if (!code || joinAttempted.value) {
      return;
    }
    joinAttempted.value = true;
    if (!isGuest.value) {
      const toast = useToast();
      toast.toast.info(t('invite.needAuth'));
      await router.replace({ name: 'room', params: { locale: locale.value } });
      return;
    }
    await ensureJoined(code);
  },
  { immediate: true },
);

watch(
  players,
  (list) => {
    if (!roomCode.value) {
      return;
    }
    if (list.length < 2 && route.name === 'room-play') {
      void router.replace({
        name: 'room-code',
        params: { locale: locale.value, code: roomCode.value },
      });
    }
  },
  { flush: 'post' },
);

function chooseGame(gameKey: string) {
  void router.push({ name: 'game', params: { locale: locale.value, game: gameKey } });
}

onBeforeRouteLeave((to) => {
  if (to.name === 'room' || to.name === 'room-code' || to.name === 'room-play') {
    return;
  }
  void leaveCurrentRoom(roomCode.value);
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
            type="button"
            class="rounded-xl border border-white/20 bg-green-500/90 px-4 py-4 text-center text-lg font-bold text-white shadow-md transition hover:bg-green-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
            @click="chooseGame(game.key)"
          >
            {{ t(`games.${game.key}`) }}
          </button>
        </div>
      </div>

      <!-- Back to Lobby -->
      <RouterLink
        :to="roomLobbyPath(locale, roomCode)"
        class="text-center text-sm font-semibold text-blue-100 underline underline-offset-4 hover:text-white"
      >
        {{ t('invite.backToLobby') }}
      </RouterLink>
    </div>
  </div>
</template>
