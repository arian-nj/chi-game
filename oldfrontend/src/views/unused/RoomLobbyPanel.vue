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
    <p v-if="isBusy" class="text-center text-sky-200/90">{{ t('invite.joining') }}</p>

    <div
      v-else
      class="flex flex-col gap-6 rounded-2xl border border-white/15 bg-custom-lite-blue/55 p-5 shadow-lg shadow-black/25 backdrop-blur-sm"
      :dir="textDir"
    >
      <div class="flex flex-col gap-2">
        <span class="text-sm font-semibold uppercase tracking-wide text-cyan-200/90">{{ t('invite.codeLabel') }}</span>
        <div class="flex items-center gap-2">
          <code
            class="flex-1 rounded-xl border border-cyan-400/30 bg-custom-deep-blue px-4 py-3 text-2xl font-bold tracking-widest text-cyan-50 shadow-inner"
          >
            {{ roomCode }}
          </code>
          <button
            type="button"
            class="rounded-xl border border-cyan-400/35 bg-cyan-500/15 px-4 py-3 text-sm font-bold text-cyan-100 transition-colors hover:bg-cyan-500/25"
            @click="copyText(roomCode)"
          >
            {{ t('invite.copy') }}
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-semibold uppercase tracking-wide text-sky-200/90">{{ t('invite.linkLabel') }}</span>
        <button
          type="button"
          class="rounded-xl border border-sky-400/25 bg-sky-500/10 px-4 py-3 text-left text-sm text-sky-100 break-all transition-colors hover:bg-sky-500/20"
          @click="copyText(roomLink)"
        >
          {{ roomLink }}
        </button>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-semibold uppercase tracking-wide text-emerald-200/90">{{ t('invite.playersLabel') }}</span>
        <p v-if="isError" class="text-sm text-red-300">{{ error?.message ?? t('invite.loadError') }}</p>
        <ul v-else class="flex flex-col gap-2">
          <li
            v-for="player in players"
            :key="String(player.id)"
            class="rounded-xl border border-white/15 bg-custom-deep-blue/75 px-4 py-3 text-lg font-semibold text-white shadow-sm"
          >
            {{ playerLabel(player.displayName, player.username) }}
          </li>
          <li
            v-if="players.length < 2"
            class="rounded-xl border border-dashed border-amber-400/35 bg-amber-500/10 px-4 py-3 text-sm text-amber-100/90"
          >
            {{ t('invite.waitingForGuest') }}
          </li>
        </ul>
      </div>

      <div class="flex flex-col gap-3 pt-1">
        <button
          type="button"
          :class="[
            'w-full rounded-xl border py-3 font-bold transition-colors',
            isReadyToPlay && isHost
              ? 'border-green-400/40 bg-green-500 text-white shadow-md shadow-green-500/30 hover:bg-green-400'
              : 'border-white/15 bg-custom-deep-blue/70 text-blue-200/70 hover:bg-custom-deep-blue/90 hover:text-blue-100',
          ]"
          @click="onStartRoom"
        >
          {{ t('invite.startRoom') }}
        </button>

        <button
          type="button"
          class="w-full rounded-xl border border-red-400/25 bg-red-950/35 py-3 font-bold text-red-100 transition-colors hover:bg-red-900/45"
          @click="leave"
        >
          {{ t('invite.leaveRoom') }}
        </button>
      </div>
    </div>
  </div>
</template>
