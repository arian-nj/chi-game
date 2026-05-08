<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';


const route = useRoute();

const navItems = [
  { path: '/profile', label: '👤', icon: 'user-o' },
  { path: '/', label: '🏠', icon: 'main-o' },
  { path: '/friends', label: '🫂', icon: 'friends-o' },
];
const showBottomNav = computed(() => {
  return route.meta.showBottomNav !== false; // Default to true if not specified
});

</script>

<template>     
    <div 
      v-if="showBottomNav"
    class="fixed bottom-0 bg-custom-blue backdrop-blur-lg
     rounded-2xl shadow-2xl border-0 flex justify-around items-center h-14 z-10 px-2
     w-11/12 sm:w-9/12 lg:w-6/12 left-1/2 -translate-x-1/2
    ">
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
          :class="isActive ? 'bg-custom-lite-blue' : 'hover:bg-gray-200/50'"
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
</template>