<script setup lang="ts">
import { RouterView } from 'vue-router';
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

  <RouterView />
  <Toast />

</template>
