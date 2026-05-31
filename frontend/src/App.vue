<script setup lang="ts">


import { onMounted } from 'vue';
import Toast, { configureToast } from './components/Toast.vue';
import { useToast } from './components/Toast.vue';

onMounted(() => {
  configureToast({
    duration: 2000,
    position: 'top-center',
    pauseOnHover: false,
    limit: 2,
  });
});

function showToast() {
  useToast().toast.success('Hello, world!');
}

</script>

<template>

  <RouterView />
  <Toast />
</template>

<style scoped></style>
