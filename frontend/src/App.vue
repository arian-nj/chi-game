<script setup lang="ts">
import { RouterView } from 'vue-router';
import Toast, { configureToast } from './components/Toast.vue';
import { onMounted } from 'vue';
import BottomNav from './components/BottomNav.vue';


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

  <div class="flex flex-col h-screen bg-gray-50 text-gray-200">
    <div class="flex-1 overflow-y-auto">
        <RouterView />
    </div>
      <BottomNav />
  </div>

  <Toast />

</template>
