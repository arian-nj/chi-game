<script setup lang="ts">
export interface TurnPlayer {
  key: string;
  label: string;
  markerClass: string;
  markerType: 'disc' | 'symbol';
  markerText?: string;
}

withDefaults(
  defineProps<{
    message: string;
    status?: 'playing' | 'thinking' | 'win' | 'draw';
    players: TurnPlayer[];
    activePlayerKey: string | null;
    showPlayers?: boolean;
  }>(),
  {
    status: 'playing',
    showPlayers: true,
  },
);

const bannerClass = {
  playing: 'border-white/20 bg-custom-lite-blue/50',
  thinking: 'border-yellow-300/50 bg-yellow-500/10',
  win: 'border-green-400/50 bg-green-500/15',
  draw: 'border-white/25 bg-custom-lite-blue/40',
};
</script>

<template>
  <div
    class="w-full rounded-2xl border p-4 shadow-md transition-colors duration-300"
    :class="bannerClass[status]"
    role="status"
    aria-live="polite"
  >
    <p
      class="text-center text-xl font-extrabold tracking-wide sm:text-2xl"
      :class="status === 'win' ? 'text-green-200' : status === 'thinking' ? 'text-yellow-200' : 'text-white'"
    >
      {{ message }}
    </p>

    <div v-if="showPlayers && players.length > 0" class="mt-4 grid grid-cols-2 gap-3">
      <div
        v-for="player in players"
        :key="player.key"
        class="flex items-center justify-center gap-3 rounded-xl border-2 px-3 py-3 transition-all duration-300"
        :class="
          player.key === activePlayerKey
            ? 'turn-indicator-active scale-[1.03] border-white bg-white/20 shadow-lg'
            : 'border-white/10 bg-black/10 opacity-45'
        "
        :aria-current="player.key === activePlayerKey ? 'true' : undefined"
      >
        <span
          v-if="player.markerType === 'disc'"
          class="size-8 shrink-0 rounded-full shadow-inner ring-2 ring-white/30"
          :class="player.markerClass"
        ></span>
        <span
          v-else
          class="flex size-8 shrink-0 items-center justify-center text-2xl font-extrabold"
          :class="player.markerClass"
        >
          {{ player.markerText }}
        </span>
        <span
          class="text-sm font-bold sm:text-base"
          :class="player.key === activePlayerKey ? 'text-white' : 'text-blue-100/70'"
        >
          {{ player.label }}
        </span>
      </div>
    </div>
  </div>
</template>
