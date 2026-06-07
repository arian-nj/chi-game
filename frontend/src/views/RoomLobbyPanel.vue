<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { useRoomSession } from '@/composables/use-room-session';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const router = useRouter();
const { toast } = useToast();

const {
  locale,
  roomCode,
  isBusy,
  players,
  isError,
  error,
  isReadyToPlay,
  isHost,
  roomLink,
  textDir,
  leave,
  playerLabel,
  copyText,
} = useRoomSession();

async function onStartRoom() {
  if (!isReadyToPlay.value) {
    toast.error(t('invite.notEnoughPlayers'));
    return;
  }
  if (!isHost.value) {
    toast.error(t('invite.notHost'));
    return;
  }
  await router.push({ name: 'room-play', params: { locale: locale.value, code: roomCode.value } });
}
</script>

<template>


  <div class="w-full pb-24 lg:pb-0 lg:p-0">
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
            {{ roomCode }}
          </code>
          <button
            type="button"
            class="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-sm font-bold text-white hover:bg-white/20"
            @click="copyText(roomCode)"
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
          @click="copyText(roomLink)"
        >
          {{ roomLink }}
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
        :class="['w-full rounded-lg border border-white/20 py-3 font-bold text-white',
         (isReadyToPlay && isHost) ? 'bg-green-500 hover:bg-green-400' : 'bg-white/10 hover:bg-white/20']"
        @click="onStartRoom"
      >
        {{ t('invite.startRoom') }}
      </button>

      <button
        type="button"
        class="w-full rounded-lg border border-white/20 bg-white/10 py-3 font-bold text-white hover:bg-white/20"
        @click="leave"
      >
        {{ t('invite.leaveRoom') }}
      </button>
    </div>
  </div>
</template>
