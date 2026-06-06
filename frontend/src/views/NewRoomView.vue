<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { useTextDirection } from '@/composables/use-text-direction';
import { createRoom, joinRoomWithCode } from '@/libs/room-api';
import { ConnectError } from '@connectrpc/connect';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';


const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const { textDir } = useTextDirection();
const { isGuest } = useGuestAuth();

const locale = computed(() => route.params.locale as string);

const joinCodeInput = ref('');
const isBusy = ref(false);

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

async function goToLobby(code: string) {
  const normalized = code.trim().toUpperCase();
  if (!normalized) {
    return;
  }
  await router.replace({
    name: 'room-code',
    params: { locale: locale.value, code: normalized },
  });
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
    await goToLobby(normalized);
  } catch (err) {
    const specific = inviteErrorMessage(err);
    toast.toast.error(specific || t('invite.joinFailed'));
  } finally {
    isBusy.value = false;
  }
}

async function onCreateRoom() {
  const toast = useToast();
  if (!isGuest.value) {
    toast.toast.info(t('invite.needAuth'));
    return;
  }

  joinCodeInput.value = '';
  isBusy.value = true;
  try {
    const created = await createRoom();
    await goToLobby(created);
  } catch (err) {
    const specific = inviteErrorMessage(err);
    toast.toast.error(specific || t('invite.createFailed'));
  } finally {
    isBusy.value = false;
  }
}

async function onJoinWithInput() {
  await joinRoom(joinCodeInput.value);
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
      {{ t('invite.roomTitle') }}
    </h1>

    <div class="w-full max-w-3xl px-4 flex-1 pb-10">
      <div
        class="bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md flex flex-col gap-5 p-5"
        :dir="textDir"
      >
        <p class="text-sm text-blue-100/80">{{ t('invite.roomIntro') }}</p>

        <button
          type="button"
          class="w-full rounded-lg bg-green-500 py-3 text-xl font-extrabold text-white disabled:opacity-60"
          :disabled="isBusy"
          @click="onCreateRoom"
        >
          {{ t('play.createRoom') }}
        </button>

        <div class="flex flex-col gap-2">
          <span class="text-sm font-semibold uppercase tracking-wide text-blue-100/90">{{ t('settings.inviteCode') }}</span>
          <input
            v-model="joinCodeInput"
            type="text"
            dir="ltr"
            autocapitalize="characters"
            spellcheck="false"
            maxlength="8"
            class="rounded-xl border border-white/20 bg-custom-deep-blue/80 px-4 py-3 text-lg font-bold tracking-widest text-white uppercase placeholder:normal-case placeholder:font-normal placeholder:tracking-normal placeholder:text-blue-100/50"
            :placeholder="t('settings.inviteCodePlaceholder')"
          />
          <button
            type="button"
            class="w-full rounded-lg border border-white/20 bg-white/10 py-3 font-bold text-white hover:bg-white/20 disabled:opacity-60"
            :disabled="isBusy || !joinCodeInput.trim()"
            @click="onJoinWithInput"
          >
            {{ t('play.joinRoom') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
