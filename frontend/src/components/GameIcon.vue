<script lang="ts" setup>
import type { GameData } from '@/lib/types';

const props = defineProps<{
    game: GameData,
}>()

</script>

<template>
    <a  :href="`/game/${props.game.key}`" class="flex flex-col items-center justify-center py-4
    font-bold bg-[#213951]/10 rounded-xl">
        <img :src="`/public/games/${props.game.key}/${props.game.key}_logo.png`" :alt="props.game.name" class="w-11 h-11">
        <p>{{ props.game.name }}</p>
    </a>
</template>

