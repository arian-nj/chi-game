<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { useTextDirection } from '@/composables/use-text-direction';
import { InviteService } from '@/gen/invite/v1/invite_pb';
import { createApiClient } from '@/libs/api-client';
import { joinRoomWithCode, leaveCurrentRoom } from '@/libs/invite-room';
import { roomInviteUrl } from '@/libs/room-url';
import { ConnectError } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { computed, ref, watch } from 'vue';
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useGuestProfile } from '@/composables/use-guest-profile';

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { textDir } = useTextDirection();
const { isGuest } = useGuestAuth();
const toast = useToast();

const locale = computed(() => route.params.locale as string);
const inviteCode = computed(() => {
  const raw = route.params.code;
  if (typeof raw !== 'string') {
    return '';
  }
  return raw.trim().toUpperCase();
});

const joinAttempted = ref(false);
const isBusy = ref(false);

const client = createApiClient(InviteService);

const { data, isError, error } = useQuery({
  queryKey: computed(() => ['invite-room', inviteCode.value]),
  queryFn: ({ signal }) => client.getInviteRoom({ inviteCode: inviteCode.value }, { signal }),
  refetchInterval: 2000,
  enabled: computed(() => Boolean(inviteCode.value)),
});
const {data:meData} = useGuestProfile();

const players = computed(() => data.value?.players ?? []);
const isReadyToPlay = computed(() => players.value.length >= 2);
const isHost = computed(() => data.value?.hostPlayer?.id === meData.value?.account?.id);

const inviteLink = computed(() => roomInviteUrl(locale.value, inviteCode.value));

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

async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text);
    const toast = useToast();
    toast.toast.success(t('invite.copied'));
  } catch {
    // ignore
  }
}

async function joinRoom(code: string) {
  const toast = useToast();
  if (!isGuest.value) {
    toast.toast.info(t('invite.needAuth'));
    return;
  }

  const normalized = code.trim().toUpperCase();
  if (!normalized) {
    return;
  }

  isBusy.value = true;
  try {
    await joinRoomWithCode(normalized);
    if (route.params.code !== normalized) {
      await router.replace({
        name: 'room-code',
        params: { locale: locale.value, code: normalized },
      });
    }
  } catch (err) {
    const specific = inviteErrorMessage(err);
    toast.toast.error(specific || t('invite.joinFailed'));
    await router.push({ name: 'room', params: { locale: locale.value } });
  } finally {
    isBusy.value = false;
  }
}

watch(
  inviteCode,
  async (code) => {
    if (!code || joinAttempted.value) {
      return;
    }
    joinAttempted.value = true;
    await joinRoom(code);
  },
  { immediate: true },
);

// watch(
//   isReadyToPlay,
//   (ready) => {
//     if (!ready || route.name !== 'room-code') {
//       return;
//     }
//     void router.replace({
//       name: 'room-play',
//       params: { locale: locale.value, code: inviteCode.value },
//     });
//   },
//   { flush: 'post' },
// );

async function onStartRoom() {
  if (!isReadyToPlay.value) {
    toast.toast.error(t('invite.notEnoughPlayers'));
    return;
  }
  if (!isHost.value) {
    toast.toast.error(t('invite.notHost'));
    return;
  }
  await router.push({ name: 'room-play', params: { locale: locale.value, code: inviteCode.value } });
}

async function onLeave() {
  try {
    await leaveCurrentRoom(inviteCode.value);
  } catch {
    // still leave UI
  }
  await router.push({ name: 'room', params: { locale: locale.value } });
}

onBeforeRouteLeave((to) => {
  if (to.name === 'room' || to.name === 'room-code' || to.name === 'room-play') {
    return;
  }
  void leaveCurrentRoom(inviteCode.value);
});
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
      {{ t('invite.roomTitle') }}
    </h1>

    <div class="w-full max-w-3xl px-4 flex-1 pb-10">
      <p v-if="isBusy" class="text-center text-blue-100/80">{{ t('invite.joining') }}</p>

      <div
        v-else
        class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md flex flex-col gap-5 p-5"
        :dir="textDir"
      >
        <div class="flex flex-col gap-2">
          <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.codeLabel') }}</span>
          <div class="flex items-center gap-2">
            <code class="flex-1 rounded-xl bg-custom-deep-blue/80 px-4 py-3 text-2xl font-bold tracking-widest text-white">
              {{ inviteCode }}
            </code>
            <button
              type="button"
              class="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-sm font-bold text-white hover:bg-white/20"
              @click="copyText(inviteCode)"
            >
              {{ t('invite.copy') }}
            </button>
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.linkLabel') }}</span>
          <button
            type="button"
            class="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-left text-sm text-blue-100 break-all hover:bg-white/20"
            @click="copyText(inviteLink)"
          >
            {{ inviteLink }}
          </button>
        </div>

        <div class="flex flex-col gap-2">
          <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('invite.playersLabel') }}</span>
          <p v-if="isError" class="text-sm text-red-300">{{ error?.message ?? t('invite.loadError') }}</p>
          <ul v-else class="flex flex-col gap-2">
            <li
              v-for="player in players"
              :key="String(player.id)"
              class="rounded-xl border border-white/10 bg-custom-deep-blue/60 px-4 py-3 text-lg font-semibold text-white"
            >
              {{ playerLabel(player.displayName, player.username) }}
            </li>
            <li
              v-if="players.length < 2"
              class="rounded-xl border border-dashed border-white/20 px-4 py-3 text-sm text-blue-100/70"
            >
              {{ t('invite.waitingForGuest') }}
            </li>
          </ul>
        </div>

        <button
          type="button"
          :class="['w-full rounded-lg border border-white/20 py-3 font-bold text-whit',
           (isReadyToPlay && isHost) ? 'bg-green-500 hover:bg-green-400' : 'bg-white/10 hover:bg-white/20']"
          @click="onStartRoom"
        >
          {{ t('invite.startRoom') }}
        </button>

        <button
          type="button"
          class="w-full rounded-lg border border-white/20 bg-white/10 py-3 font-bold text-white hover:bg-white/20"
          @click="onLeave"
        >
          {{ t('invite.leaveRoom') }}
        </button>
      </div>
    </div>
  </div>
</template>
