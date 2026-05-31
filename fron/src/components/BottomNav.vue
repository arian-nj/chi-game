<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const navItems = [
  { path: '/me', label: 'پروفایل', icon: '👤' },
  { path: '/', label: 'خانه', icon: '🏠' },
  { path: '/inbox', label: 'چت', icon: '💬' },
];

const showBottomNav = computed(() => route.meta.showBottomNav !== false);
</script>

<template>
  <transition name="fade" mode="out-in">
    <div
      v-if="showBottomNav"
      class="fixed bottom-4 left-1/2 -translate-x-1/2 w-11/12 sm:w-9/12 lg:w-6/12 z-30
             backdrop-blur-2xl border-0 rounded-3xl shadow-[0_12px_40px_0px_rgba(0,71,255,0.15)] 
             flex justify-around items-center h-20 px-2
             ring-2 ring-blue-400/20 animate-bounceInUp"
      style="box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.30);"
    >
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        custom
        v-slot="{ navigate, isActive }"
      >
        <button
          @click="navigate"
          class="flex flex-col items-center justify-center flex-1 h-full mx-2 py-2 transition-all duration-200
                 rounded-xl group relative"
          :class="isActive 
                    ? 'bg-white/30 shadow-md scale-105'
                    : 'hover:bg-white/10 hover:scale-105'"
        >
          <span
            class="text-4xl transition-all duration-200 drop-shadow-md"
            :class="isActive ? 'text-blue-700 scale-125' : 'text-blue-100 group-hover:text-blue-200'"
          >
            {{ item.icon }}
          </span>
          <span
            class="text-xs font-semibold mt-1 transition-all duration-200 tracking-wide"
            :class="isActive ? 'text-blue-900' : 'text-blue-50 group-hover:text-blue-200 opacity-70'"
          >
            {{ item.label }}
          </span>
          <span
            v-if="isActive"
            class="absolute bottom-1 left-1/2 -translate-x-1/2 w-6 h-1 rounded-xl bg-blue-500/90 animate-pulse"
          ></span>
        </button>
      </router-link>
    </div>
  </transition>
</template>

<style scoped>
@keyframes bounceInUp {
  0% {
    opacity: 0;
    transform: translateY(80px) scale(0.97);
  }
  70% {
    opacity: 1;
    transform: translateY(-6px) scale(1.01);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
.animate-bounceInUp {
  animation: bounceInUp 0.70s cubic-bezier(.41,2.41,.52,.97) both;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s cubic-bezier(.4,2,.75,.82);
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>