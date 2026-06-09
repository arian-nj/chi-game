<script setup lang="ts">
import type { BotDifficulty } from '@/libs/bot-difficulty';
import { botDifficulties } from '@/libs/bot-difficulty';

const difficulty = defineModel<BotDifficulty>({ default: 'medium' });

defineProps<{
  labels: Record<BotDifficulty, string>;
}>();
</script>

<template>
  <div
    class="grid grid-cols-3 gap-1.5 rounded-2xl border border-white/10 bg-custom-deep-blue/90 p-1"
    role="group"
  >
    <button
      v-for="level in botDifficulties"
      :key="level"
      type="button"
      class="rounded-xl px-2 py-2.5 text-sm font-bold transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50"
      :class="
        difficulty === level
          ? 'bg-white/90 text-custom-blue shadow'
          : 'text-blue-100/80 hover:bg-white/10 hover:text-white'
      "
      :aria-pressed="difficulty === level"
      @click="difficulty = level"
    >
      {{ labels[level] }}
    </button>
  </div>
</template>
