<script setup lang="ts">
import { useGuestAuth } from '@/composables/use-guest-auth';
import { usePageSeo } from '@/composables/use-page-seo';
import Toast, { configureToast } from './components/Toast.vue';
import { onMounted } from 'vue';

usePageSeo();
useGuestAuth();

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
    <RouterView />
    <Toast />

  </div>
</template>
