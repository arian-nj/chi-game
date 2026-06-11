<script lang="ts" setup>
import type { GameData } from '@/lib/types';
import { useLocalePath } from '@/composables/useLocalePath';

const props = defineProps<{
    game: GameData,
}>()

const { localePath } = useLocalePath();

</script>

<template>
    <RouterLink :to="localePath(`/game/${props.game.key}`)" class="flex flex-col items-center justify-center py-4
    font-bold bg-[#213951]/10 rounded-xl">
        <img :src="props.game.logo" :alt="props.game.name" class="w-11 h-11">
        <p>{{ props.game.name }}</p>
    </RouterLink>
</template>

