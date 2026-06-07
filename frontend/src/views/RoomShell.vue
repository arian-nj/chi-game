<script setup lang="ts">
import RoomChat from '@/components/RoomChat.vue';
import { provideRoomSession, roomSessionKey } from '@/composables/use-room-session';
import { provide } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const session = provideRoomSession();
provide(roomSessionKey, session);
</script>

<template>
  <div class="relative min-h-screen bg-custom-blue text-white">
    <RouterLink
      :to="`/${session.locale}`"
      class="absolute left-4 top-4 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :aria-label="t('nav.backToHome')"
    >
      <span aria-hidden="true">←</span>
      {{ t('nav.home') }}
    </RouterLink>

    <div class="flex flex-col pt-14 lg:grid lg:grid-cols-[1fr_minmax(0,42rem)_1fr] lg:items-start lg:gap-x-6 lg:px-6">
      <div class="hidden lg:block" aria-hidden="true"></div>

      <main class="w-full px-4 pb-24 lg:px-0 lg:pb-10">
        <RouterView />
      </main>

      <div class="sticky top-4 w-80 justify-self-end">
        <RoomChat />
      </div>
    </div>
  </div>
</template>
