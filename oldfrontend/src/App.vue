<script setup lang="ts">
import { usePageSeo } from '@/composables/use-page-seo';
import Toast, { configureToast } from './components/Toast.vue';
import { onMounted } from 'vue';

usePageSeo();

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
