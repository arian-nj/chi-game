<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { gamesData } from '../libs/game';

const route = useRoute();
const urlGameName = route.params.game;
const gameData = gamesData.find(game => game.key === urlGameName);

// Tabs
const tabs = [
    { name: 'Play', icon: '🚀', to: `/game/${urlGameName}` },
    { name: 'Rules', icon: '📖', to: `/game/${urlGameName}/rules` },
];

// Play Button Animation
const playBtnRef = ref<HTMLButtonElement | null>(null);

    onMounted(() => {
  if (playBtnRef.value) {
    playBtnRef.value.style.transform = 'translateY(120%)';
    playBtnRef.value.style.opacity = '0';
    setTimeout(() => {
      playBtnRef.value!.style.transition = 'transform 0.7s cubic-bezier(.23,1.12,.75,.95), opacity 0.7s cubic-bezier(.23,1.12,.75,.95)';
      playBtnRef.value!.style.transform = 'translateY(0%)';
      playBtnRef.value!.style.opacity = '1';
    }, 30);
  }
});

function playGame() {
    const toast = useToast();
    toast.toast.success('Game started');
}

</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center pb-3 text-white relative">
    <h1 class="text-4xl font-bold text-white mb-8 mt-2 animate-pop select-none drop-shadow-sm uppercase">
      {{ gameData?.name }}
    </h1>
    
    <!-- Tab Navigation -->
    <div class="flex space-x-4 bg-custom-lite-blue/70 rounded-xl mb-6 px-2 py-1 shadow-lg">
        <RouterLink
            v-for="tab in tabs"
            :key="tab.name"
            :to="tab.to"
            custom
            v-slot="{ isExactActive, navigate }"
        >

            <a
            :href="tab.to"
            @click="navigate"
            :class="[
                'px-6 py-2 rounded-lg transition font-bold flex items-center gap-2 cursor-pointer',
                isExactActive /* USING IT HERE NOW */
                ? 'bg-white/90 text-custom-blue shadow'
                : 'bg-transparent text-blue-100 hover:bg-white/30'
            ]"
            >
            <span>{{ tab.icon }}</span> {{ tab.name }}
            </a>
        </RouterLink>
    </div>

    <div class="w-full max-w-3xl px-4 flex-1">
      <RouterView />
    </div>

    <!-- Play Button (Only show if on the main play tab, or keep visible everywhere based on your design) -->
    <div class="flex flex-col items-center justify-center">
        <button
          ref="playBtnRef"
          class="bg-green-500 text-white p-3 rounded-lg w-3/4 mx-auto text-3xl font-extrabold fixed left-1/2 -translate-x-1/2 bottom-10 z-20"
          style="opacity: 0; transform: translateY(120%);"
          @click="playGame"
        >
          <span class="text-3xl animate-bounce inline-block">🚀</span><span>Play</span>
        </button>
    </div>
  </div>
</template>