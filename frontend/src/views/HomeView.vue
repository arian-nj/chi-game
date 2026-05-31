<script setup lang="ts">
import { useRouter } from "vue-router";

const router = useRouter();

const games = [
  { key: "xo3x3", name: "XO 3x3" },
  { key: "conn4", name: "Connect 4" },
  { key: "linedot", name: "Line Dot" },
  { key: "go", name: "Go!" }
];

function goToGame(key: string) {
  router.push(`/game/${key}`);
}
</script>

<template>
  <div class="bg-custom-blue min-h-screen w-screen flex flex-col items-center p-6">

    <h1 class="text-4xl font-bold text-white mb-8 mt-2 animate-pop select-none drop-shadow-sm">
      Chi Game
    </h1>

    <div 
      class="grid w-full max-w-5xl gap-6
             grid-cols-2 
             md:grid-cols-3 
             lg:grid-cols-4
             place-items-center
             mx-auto
             px-4
      "
    >
      <button
        v-for="game in games"
        :key="game.key"
        @click="goToGame(game.key)"
        class="
          bg-custom-lite-blue hover:bg-custom-deep-blue 
          rounded-2xl 
          w-full
          transition 
          duration-100 
          shadow 
          outline-none
          focus:ring-2 focus:ring-blue-500
          aspect-4/3
          overflow-hidden
          flex flex-col items-center justify-center
        "
      >
        <div class="flex flex-col items-center justify-center w-full h-full p-4">
          <span class="text-lg md:text-xl font-semibold text-white select-none text-center">
            {{ game.name }}
          </span>
        </div>
      </button>
    </div>
  </div>
</template>