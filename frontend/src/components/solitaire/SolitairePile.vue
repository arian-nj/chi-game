<script setup lang="ts">
import { computed } from 'vue';
import type { Card, DrawMode, PileRef } from '@/lib/solitaire/types';
import SolitaireCard from './SolitaireCard.vue';

const props = defineProps<{
    cards: Card[];
    variant: 'stock' | 'waste' | 'foundation' | 'tableau';
    pileRef: PileRef | { kind: 'stock' };
    selectedCardIndex?: number | null;
    highlighted?: boolean;
    stockHasCards?: boolean;
    wasteFanCount?: DrawMode;
    draggingFromIndex?: number | null;
}>();

const emit = defineEmits<{
    cardClick: [cardIndex: number];
    cardPointerDown: [cardIndex: number, event: PointerEvent];
    slotClick: [];
}>();

const pileDataIndex = computed(() => {
    if (props.pileRef.kind === 'foundation' || props.pileRef.kind === 'tableau') {
        return props.pileRef.index;
    }
    return undefined;
});

function isDragHidden(index: number): boolean {
    return (
        props.draggingFromIndex !== null &&
        props.draggingFromIndex !== undefined &&
        index >= props.draggingFromIndex
    );
}

const visibleWasteCards = computed(() => {
    if (props.variant !== 'waste') {
        return [];
    }
    const fanCount = props.wasteFanCount ?? 3;
    const start = Math.max(0, props.cards.length - fanCount);
    return props.cards.slice(start).map((card, offset) => ({
        card,
        index: start + offset,
    }));
});

const pileHeight = computed(() => {
    if (props.variant === 'tableau') {
        const count = Math.max(props.cards.length, 1);
        return `calc(var(--card-height) + (var(--card-overlap) * ${count - 1}))`;
    }
    return 'var(--card-height)';
});

const wasteWidth = computed(() => {
    const fanCount = props.wasteFanCount ?? 3;
    const visible = Math.min(props.cards.length, fanCount);
    if (visible <= 1) {
        return 'var(--card-width)';
    }
    return `calc(var(--card-width) + (var(--card-fan) * ${visible - 1}))`;
});
</script>

<template>
    <div
        class="sol-pile"
        :class="[
            `sol-pile-${variant}`,
            { 'sol-pile-highlight': highlighted },
        ]"
        :data-pile-kind="pileRef.kind"
        :data-pile-index="pileDataIndex"
        :style="{
            minHeight: pileHeight,
            width: variant === 'waste' ? wasteWidth : 'var(--card-width)',
        }"
    >
        <button
            v-if="variant === 'stock'"
            type="button"
            class="sol-slot sol-stock-slot"
            :class="{ 'sol-stock-slot-active': stockHasCards }"
            @click="emit('slotClick')"
        >
            <span v-if="stockHasCards" class="sol-stock-back" aria-hidden="true" />
        </button>

        <button
            v-else-if="cards.length === 0"
            type="button"
            class="sol-slot"
            @click="emit('slotClick')"
        />

        <template v-else-if="variant === 'waste'">
            <div
                v-for="(entry, offset) in visibleWasteCards"
                :key="entry.index"
                class="sol-pile-card-wrap sol-pile-card-fan"
                :style="{
                    left: `calc(var(--card-fan) * ${offset})`,
                    zIndex: offset,
                }"
            >
                <SolitaireCard
                    :card="entry.card"
                    :selected="selectedCardIndex === entry.index"
                    :drag-hidden="isDragHidden(entry.index)"
                    @click="emit('cardClick', entry.index)"
                    @pointer-down="emit('cardPointerDown', entry.index, $event)"
                />
            </div>
        </template>

        <template v-else-if="variant === 'foundation'">
            <div class="sol-pile-card-wrap">
                <SolitaireCard
                    :card="cards[cards.length - 1]!"
                    :selected="selectedCardIndex === cards.length - 1"
                    :drag-hidden="isDragHidden(cards.length - 1)"
                    @click="emit('cardClick', cards.length - 1)"
                    @pointer-down="emit('cardPointerDown', cards.length - 1, $event)"
                />
            </div>
        </template>

        <template v-else>
            <div
                v-for="(card, index) in cards"
                :key="index"
                class="sol-pile-card-wrap"
                :style="{ top: `calc(var(--card-overlap) * ${index})` }"
            >
                <SolitaireCard
                    :card="card"
                    :selected="
                        selectedCardIndex !== null &&
                        selectedCardIndex !== undefined &&
                        index >= selectedCardIndex
                    "
                    :drag-hidden="isDragHidden(index)"
                    @click="emit('cardClick', index)"
                    @pointer-down="emit('cardPointerDown', index, $event)"
                />
            </div>
        </template>
    </div>
</template>

<style scoped>
.sol-pile {
    position: relative;
    flex-shrink: 0;
}

.sol-pile-highlight .sol-slot,
.sol-pile-highlight {
    box-shadow: 0 0 0 2px #2e8b2e, 0 0 10px rgba(46, 139, 46, 0.45);
    border-radius: 4px;
}

.sol-slot {
    width: var(--card-width);
    height: var(--card-height);
    padding: 0;
    border-radius: 4px;
    border: 2px dashed #6b8f6b;
    background: rgba(0, 80, 0, 0.08);
    box-sizing: border-box;
}

.sol-stock-slot {
    border-style: solid;
    border-color: #808080;
    background: #c0c0c0;
}

.sol-stock-slot-active {
    border-style: solid;
    border-color: #4a4a4a;
    background: transparent;
}

.sol-stock-back {
    display: block;
    width: 100%;
    height: 100%;
    border-radius: 3px;
    background:
        repeating-linear-gradient(
            45deg,
            #1e4f9a 0,
            #1e4f9a 4px,
            #2a63b8 4px,
            #2a63b8 8px
        );
    border: 1px solid #0f2f66;
    box-sizing: border-box;
}

.sol-pile-card-wrap {
    position: absolute;
    left: 0;
    top: 0;
}

.sol-pile-card-fan {
    top: 0;
}
</style>
