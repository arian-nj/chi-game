<script setup lang="ts">
import { provideRoomSession, roomSessionKey } from '@/composables/use-room-session';
import { provide } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const session = provideRoomSession();
provide(roomSessionKey, session);
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center pb-3 pt-14 text-white relative">
    <RouterLink
      :to="`/${session.locale}`"
      class="absolute left-4 top-14 z-30 flex items-center gap-2 rounded-xl border border-white/10 bg-custom-lite-blue/70 px-4 py-2 text-sm font-bold text-blue-100 shadow-md transition hover:bg-white/20 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :aria-label="t('nav.backToHome')"
    >
      <span aria-hidden="true">←</span>
      {{ t('nav.home') }}
    </RouterLink>

    <RouterView />
  </div>
</template>
