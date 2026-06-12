<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import SolitairePile from '@/components/solitaire/SolitairePile.vue';
import {
    applyMove,
    canSelectSource,
    flipStock,
    getValidDestinations,
    isWon,
    newGame,
    pileRefsEqual,
} from '@/lib/solitaire/solitaire';
import {
    FOUNDATION_PILE_COUNT,
    TABLEAU_PILE_COUNT,
    type GameState,
    type MoveDestination,
    type MoveSource,
    type PileRef,
} from '@/lib/solitaire/types';

const props = defineProps<{
    drawMode: 1 | 3;
}>();

const gameState = ref<GameState>(newGame(props.drawMode));
const selection = ref<MoveSource | null>(null);
const validDestinations = ref<MoveDestination[]>([]);
const elapsedSeconds = ref(0);
const facePressed = ref(false);
let timerId: ReturnType<typeof setInterval> | null = null;

const gameWon = computed(() => isWon(gameState.value));

const faceEmoji = computed(() => {
    if (facePressed.value) return '😮';
    if (gameWon.value) return '😎';
    return '🙂';
});

function formatCounter(value: number): string {
    const clamped = Math.max(0, Math.min(999, value));
    return String(clamped).padStart(3, '0');
}

function startTimer() {
    if (timerId) return;
    timerId = setInterval(() => {
        if (elapsedSeconds.value < 999) {
            elapsedSeconds.value++;
        }
    }, 1000);
}

function stopTimer() {
    if (timerId) {
        clearInterval(timerId);
        timerId = null;
    }
}

function clearSelection() {
    selection.value = null;
    validDestinations.value = [];
}

function isHighlighted(pile: PileRef): boolean {
    return validDestinations.value.some((destination) => pileRefsEqual(destination.pile, pile));
}

function getSelectedCardIndex(pile: PileRef): number | null {
    if (!selection.value || !pileRefsEqual(selection.value.pile, pile)) {
        return null;
    }
    return selection.value.cardIndex;
}

function selectCard(pile: PileRef, cardIndex: number) {
    const source: MoveSource = { pile, cardIndex };
    if (!canSelectSource(gameState.value, source)) {
        return;
    }

    selection.value = source;
    validDestinations.value = getValidDestinations(gameState.value, source);
}

function tryMoveTo(pile: PileRef) {
    if (!selection.value) {
        return;
    }

    const next = applyMove(gameState.value, selection.value, { pile });
    if (!next) {
        return;
    }

    startTimer();
    gameState.value = next;
    clearSelection();
}

function onCardClick(pile: PileRef, cardIndex: number) {
    if (gameWon.value) {
        return;
    }

    if (selection.value) {
        if (pileRefsEqual(selection.value.pile, pile) && selection.value.cardIndex === cardIndex) {
            clearSelection();
            return;
        }

        if (isHighlighted(pile)) {
            tryMoveTo(pile);
            return;
        }

        if (canSelectSource(gameState.value, { pile, cardIndex })) {
            selectCard(pile, cardIndex);
            return;
        }

        clearSelection();
        return;
    }

    selectCard(pile, cardIndex);
}

function onSlotClick(pile: PileRef | { kind: 'stock' }) {
    if (pile.kind === 'stock') {
        if (gameWon.value) {
            return;
        }
        clearSelection();
        const next = flipStock(gameState.value);
        if (next !== gameState.value) {
            startTimer();
            gameState.value = next;
        }
        return;
    }

    if (selection.value && isHighlighted(pile)) {
        tryMoveTo(pile);
    }
}

function resetGame() {
    stopTimer();
    gameState.value = newGame(props.drawMode);
    elapsedSeconds.value = 0;
    facePressed.value = false;
    clearSelection();
}

watch(
    () => props.drawMode,
    () => {
        resetGame();
    },
);

onMounted(() => {
    resetGame();
});

onUnmounted(() => {
    stopTimer();
});

defineExpose({ resetGame });
</script>

<template>
    <div
        class="xp-game-body"
        :class="{ 'xp-game-won': gameWon }"
        :style="{
            '--card-width': '52px',
            '--card-height': '72px',
            '--card-overlap': '18px',
            '--card-tight-overlap': '4px',
            '--card-fan': '18px',
            '--ui-cell-size': '26px',
        }"
    >
        <div class="xp-sunken xp-status-panel">
            <div class="xp-led">{{ formatCounter(gameState.moves) }}</div>
            <button
                type="button"
                class="xp-face-btn"
                :class="{ 'xp-face-btn-won': gameWon }"
                @click="resetGame"
                @mousedown="facePressed = true"
                @mouseup="facePressed = false"
                @mouseleave="facePressed = false"
            >{{ faceEmoji }}</button>
            <div class="xp-led">{{ formatCounter(elapsedSeconds) }}</div>
        </div>

        <div class="sol-board-area">
            <div v-if="gameWon" class="xp-win-banner" role="status" aria-live="polite">
                <span class="xp-win-title">You Win!</span>
                <span class="xp-win-time">Finished in {{ formatCounter(elapsedSeconds) }}s</span>
            </div>

            <div class="sol-board" :class="{ 'sol-board-won': gameWon }">
                <div class="sol-top-row">
                    <div class="sol-top-left">
                        <SolitairePile
                            variant="stock"
                            :cards="[]"
                            :pile-ref="{ kind: 'stock' }"
                            :stock-has-cards="gameState.stock.length > 0"
                            @slot-click="onSlotClick({ kind: 'stock' })"
                        />
                        <SolitairePile
                            variant="waste"
                            :cards="gameState.waste"
                            :pile-ref="{ kind: 'waste' }"
                            :selected-card-index="getSelectedCardIndex({ kind: 'waste' })"
                            :highlighted="isHighlighted({ kind: 'waste' })"
                            :waste-fan-count="gameState.drawMode"
                            @card-click="onCardClick({ kind: 'waste' }, $event)"
                            @slot-click="onSlotClick({ kind: 'waste' })"
                        />
                    </div>

                    <div class="sol-foundations">
                        <SolitairePile
                            v-for="index in FOUNDATION_PILE_COUNT"
                            :key="`foundation-${index - 1}`"
                            variant="foundation"
                            :cards="gameState.foundations[index - 1]!"
                            :pile-ref="{ kind: 'foundation', index: index - 1 }"
                            :selected-card-index="
                                getSelectedCardIndex({ kind: 'foundation', index: index - 1 })
                            "
                            :highlighted="isHighlighted({ kind: 'foundation', index: index - 1 })"
                            @card-click="
                                onCardClick({ kind: 'foundation', index: index - 1 }, $event)
                            "
                            @slot-click="onSlotClick({ kind: 'foundation', index: index - 1 })"
                        />
                    </div>
                </div>

                <div class="sol-tableau">
                    <SolitairePile
                        v-for="index in TABLEAU_PILE_COUNT"
                        :key="`tableau-${index - 1}`"
                        variant="tableau"
                        :cards="gameState.tableau[index - 1]!"
                        :pile-ref="{ kind: 'tableau', index: index - 1 }"
                        :selected-card-index="
                            getSelectedCardIndex({ kind: 'tableau', index: index - 1 })
                        "
                        :highlighted="isHighlighted({ kind: 'tableau', index: index - 1 })"
                        @card-click="onCardClick({ kind: 'tableau', index: index - 1 }, $event)"
                        @slot-click="onSlotClick({ kind: 'tableau', index: index - 1 })"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.xp-game-body {
    width: fit-content;
    max-width: 100%;
    margin: auto;
    padding: 8px;
    background: #0a6b0a;
}

.xp-sunken {
    border: 3px solid;
    border-color: #808080 #fff #fff #808080;
    background: #c0c0c0;
}

.xp-status-panel {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    box-sizing: border-box;
    padding: calc(var(--ui-cell-size) * 0.2) calc(var(--ui-cell-size) * 0.35);
    margin-bottom: calc(var(--ui-cell-size) * 0.3);
}

.xp-led {
    background: #000;
    color: #f00;
    font-family: 'Courier New', Courier, monospace;
    font-size: calc(var(--ui-cell-size) * 1.1);
    font-weight: bold;
    line-height: 1;
    padding: 2px 4px;
    min-width: calc(var(--ui-cell-size) * 2);
    text-align: right;
    letter-spacing: 1px;
}

.xp-face-btn {
    width: calc(var(--ui-cell-size) * 1.5);
    height: calc(var(--ui-cell-size) * 1.5);
    font-size: calc(var(--ui-cell-size) * 0.85);
    line-height: 1;
    padding: 0;
    cursor: default;
    background: #c0c0c0;
    border: 3px solid;
    border-color: #fff #808080 #808080 #fff;
    display: flex;
    align-items: center;
    justify-content: center;
}

.xp-face-btn:active {
    border-color: #808080 #fff #fff #808080;
}

.xp-face-btn-won {
    background: #fff9b0;
    animation: xp-win-face 0.55s ease-in-out infinite alternate;
}

@keyframes xp-win-face {
    from {
        transform: scale(1);
    }
    to {
        transform: scale(1.12);
    }
}

.sol-board-area {
    position: relative;
    overflow-x: auto;
    max-width: 100%;
}

.sol-board {
    padding: 8px;
    border-radius: 4px;
    transition: box-shadow 0.3s ease;
}

.sol-board-won {
    box-shadow: 0 0 0 3px #2e8b2e, 0 0 18px rgba(80, 220, 80, 0.65);
}

.sol-top-row {
    display: flex;
    justify-content: space-between;
    gap: 16px;
    margin-bottom: 20px;
}

.sol-top-left {
    display: flex;
    gap: 12px;
}

.sol-foundations {
    display: flex;
    gap: 12px;
}

.sol-tableau {
    display: flex;
    gap: 12px;
    align-items: flex-start;
}

.xp-win-banner {
    position: absolute;
    inset: 0;
    z-index: 5;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6px;
    pointer-events: none;
    background: rgba(0, 80, 0, 0.28);
}

.xp-win-title {
    padding: 8px 20px;
    background: linear-gradient(180deg, #fff9b0 0%, #ffd700 55%, #e6b800 100%);
    border: 3px solid;
    border-color: #fff #8b7500 #8b7500 #fff;
    color: #1a4d1a;
    font-size: 1.5rem;
    font-weight: bold;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
}

.xp-win-time {
    padding: 4px 12px;
    background: rgba(0, 0, 0, 0.72);
    color: #7fff7f;
    font-family: 'Courier New', Courier, monospace;
    font-size: 0.95rem;
    font-weight: bold;
}

@media (max-width: 640px) {
    .xp-game-body {
        width: 100%;
        --card-width: 44px;
        --card-height: 62px;
        --card-overlap: 15px;
        --card-fan: 14px;
    }

    .sol-top-row,
    .sol-tableau {
        gap: 8px;
    }

    .sol-top-left,
    .sol-foundations {
        gap: 8px;
    }
}
</style>
