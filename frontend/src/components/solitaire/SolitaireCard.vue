<script setup lang="ts">
import { computed } from 'vue';
import { cardLabel, suitSymbol } from '@/lib/solitaire/solitaire';
import { isRed, type Card } from '@/lib/solitaire/types';

const props = defineProps<{
    card: Card;
    selected?: boolean;
}>();

const emit = defineEmits<{
    click: [];
}>();

const colorClass = computed(() => {
    if (!props.card.faceUp) {
        return '';
    }
    return isRed(props.card.suit) ? 'sol-card-red' : 'sol-card-black';
});
</script>

<template>
    <button
        type="button"
        class="sol-card"
        :class="[
            colorClass,
            {
                'sol-card-selected': selected,
                'sol-card-back': !card.faceUp,
            },
        ]"
        @click.stop="emit('click')"
    >
        <template v-if="card.faceUp">
            <span class="sol-card-rank">{{ cardLabel(card.rank) }}</span>
            <span class="sol-card-suit">{{ suitSymbol(card.suit) }}</span>
        </template>
    </button>
</template>

<style scoped>
.sol-card {
    width: var(--card-width);
    height: var(--card-height);
    padding: 0;
    border-radius: 4px;
    border: 1px solid #4a4a4a;
    background: #fff;
    box-sizing: border-box;
    position: relative;
    cursor: default;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
    touch-action: manipulation;
}

.sol-card-back {
    background:
        repeating-linear-gradient(
            45deg,
            #1e4f9a 0,
            #1e4f9a 4px,
            #2a63b8 4px,
            #2a63b8 8px
        );
    border-color: #0f2f66;
}

.sol-card-selected {
    box-shadow: 0 0 0 2px #0054e3, 0 0 8px rgba(0, 84, 227, 0.55);
}

.sol-card-red {
    color: #c00;
}

.sol-card-black {
    color: #111;
}

.sol-card-rank {
    position: absolute;
    top: 3px;
    left: 4px;
    font-size: calc(var(--card-width) * 0.24);
    font-weight: bold;
    line-height: 1;
}

.sol-card-suit {
    position: absolute;
    left: 50%;
    top: 52%;
    transform: translate(-50%, -50%);
    font-size: calc(var(--card-width) * 0.42);
    line-height: 1;
}
</style>
