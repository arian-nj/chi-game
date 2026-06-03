<script setup lang="ts">
import LocaleSwitcher from '@/components/LocaleSwitcher.vue';
import { useBackendHealth } from '@/composables/use-backend-health';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { usePageSeo } from '@/composables/use-page-seo';
import Toast, { configureToast } from './components/Toast.vue';
import { onMounted } from 'vue';

usePageSeo();
useGuestAuth();
const { isBackendHealthy, backendError, isLoading } = useBackendHealth();
const isDev = import.meta.env.DEV;

onMounted(() => {
  configureToast({
    duration: 2000,
    position: 'top-center',
    pauseOnHover: false,
    limit: 2,
  });
});
</script>

<template>
  <div class="relative min-h-screen">
    <div class="fixed right-4 top-4 z-40">
      <LocaleSwitcher />
    </div>
    <RouterView />
    <Toast />
    <p
      v-if="isDev"
      class="fixed bottom-3 left-3 z-40 rounded-md bg-black/50 px-2 py-1 text-xs text-white/90"
      :title="backendError ?? undefined"
    >
      API:
      <span v-if="isLoading" class="text-amber-200">…</span>
      <span v-else-if="isBackendHealthy" class="text-emerald-300">online</span>
      <span v-else class="text-rose-300">offline</span>
    </p>
  </div>
</template>
