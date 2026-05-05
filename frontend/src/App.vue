<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router';
import Toast, { configureToast } from './components/Toast.vue';
import { computed, onMounted } from 'vue';

const route = useRoute();

onMounted(() => {
  configureToast({
    duration: 2000,
    position: 'top-center',
    pauseOnHover: false,
    limit: 2,
  });
});

const navItems = [
  { path: '/profile', label: '👤', icon: 'user-o' },
  { path: '/', label: '🏠', icon: 'main-o' },
  { path: '/rooms', label: '🛋️', icon: 'apps-o' },
];
const showBottomNav = computed(() => {
  return route.meta.showBottomNav !== false; // Default to true if not specified
});
</script>

<template>
  <div class="flex flex-col h-screen bg-gray-50">
    <!-- Main content -->
    <div class="flex-1 overflow-y-auto pb-14">
      <RouterView class="h-full" />
    </div>
    <Toast />

    <!-- Bottom navigation -->
     
    <div 
      v-if="showBottomNav"
    class="fixed bottom-0 left-0 right-0 bg-indigo-200 backdrop-blur-lg rounded-t-2xl shadow-2xl border-0 flex justify-around items-center h-14 z-10 px-2">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        custom
        v-slot="{ navigate, isActive }"
      >
        <button
          @click="navigate"
          class="flex flex-col items-center justify-center flex-1 h-full rounded-xl transition-all duration-200 active:scale-95 focus:outline-none"
          :class="isActive ? 'bg-gray-100/70' : 'hover:bg-gray-200/50'"
        >
          <van-icon
            :name="item.icon"
            class="text-2xl transition-transform duration-200"
            :class="isActive ? 'scale-105' : ''"
          />
          <span 
            class="text-3xl mt-0.5 transition-transform duration-200"
            :class="isActive ? 'scale-110' : ''"
          >
            {{ item.label }}
          </span>
        </button>
      </router-link>
    </div>
  </div>
</template>