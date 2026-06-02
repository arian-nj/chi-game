<script setup lang="ts">
import LocaleSwitcher from '@/components/LocaleSwitcher.vue';
import Toast, { configureToast } from './components/Toast.vue';
import { onMounted } from 'vue';

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
  </div>
</template>
