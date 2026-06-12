import {
    canPlaceOnFoundation,
    canPlaceOnTableau,
    getMovableStack,
} from './moves';
import {
    FOUNDATION_PILE_COUNT,
    RANKS,
    SUITS,
    TABLEAU_PILE_COUNT,
    type Card,
    type DrawMode,
    type GameState,
    type MoveDestination,
    type MoveSource,
    type PileRef,
    type Rank,
    type Suit,
} from './types';

export function createDeck(): Card[] {
    const deck: Card[] = [];
    for (const suit of SUITS) {
        for (const rank of RANKS) {
            deck.push({ suit, rank, faceUp: false });
        }
    }
    return deck;
}

function shuffle<T>(items: T[]): T[] {
    const arr = [...items];
    for (let i = arr.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [arr[i], arr[j]] = [arr[j]!, arr[i]!];
    }
    return arr;
}

export function createShuffledDeck(): Card[] {
    return shuffle(createDeck());
}

export function dealKlondike(drawMode: DrawMode, deck: Card[] = createShuffledDeck()): GameState {
    if (deck.length !== 52) {
        throw new Error('Klondike deal requires a 52-card deck');
    }

    const cards = deck.map((card) => ({ ...card }));
    const tableau: Card[][] = Array.from({ length: TABLEAU_PILE_COUNT }, () => []);
    let index = 0;

    for (let col = 0; col < TABLEAU_PILE_COUNT; col++) {
        for (let row = 0; row <= col; row++) {
            const card = cards[index++]!;
            tableau[col]!.push({
                suit: card.suit,
                rank: card.rank,
                faceUp: row === col,
            });
        }
    }

    const stock = cards.slice(index).map((card) => ({
        suit: card.suit,
        rank: card.rank,
        faceUp: false,
    }));

    return {
        stock,
        waste: [],
        foundations: Array.from({ length: FOUNDATION_PILE_COUNT }, () => []),
        tableau,
        drawMode,
        moves: 0,
    };
}

export function newGame(drawMode: DrawMode = 3): GameState {
    return dealKlondike(drawMode);
}

function cloneCards(cards: Card[]): Card[] {
    return cards.map((card) => ({ ...card }));
}

export function cloneGameState(state: GameState): GameState {
    return {
        stock: cloneCards(state.stock),
        waste: cloneCards(state.waste),
        foundations: state.foundations.map(cloneCards),
        tableau: state.tableau.map(cloneCards),
        drawMode: state.drawMode,
        moves: state.moves,
    };
}

function getPile(state: GameState, ref: PileRef): Card[] {
    switch (ref.kind) {
        case 'waste':
            return state.waste;
        case 'foundation':
            return state.foundations[ref.index]!;
        case 'tableau':
            return state.tableau[ref.index]!;
    }
}

function isSamePile(left: PileRef, right: PileRef): boolean {
    if (left.kind !== right.kind) {
        return false;
    }
    if (left.kind === 'waste') {
        return true;
    }
    if (left.kind === 'foundation' && right.kind === 'foundation') {
        return left.index === right.index;
    }
    if (left.kind === 'tableau' && right.kind === 'tableau') {
        return left.index === right.index;
    }
    return false;
}

function revealTableauTop(pile: Card[]): void {
    if (pile.length === 0) {
        return;
    }

    const top = pile[pile.length - 1]!;
    if (!top.faceUp) {
        pile[pile.length - 1] = { ...top, faceUp: true };
    }
}

function getMoveStack(state: GameState, from: MoveSource): Card[] | null {
    const pile = getPile(state, from.pile);
    if (pile.length === 0) {
        return null;
    }

    if (from.pile.kind === 'waste' || from.pile.kind === 'foundation') {
        const topIndex = pile.length - 1;
        if (from.cardIndex !== topIndex) {
            return null;
        }
        const topCard = pile[topIndex]!;
        return topCard.faceUp ? [topCard] : null;
    }

    return getMovableStack(pile, from.cardIndex);
}

export function flipStock(state: GameState): GameState {
    const next = cloneGameState(state);

    if (next.stock.length > 0) {
        const count = Math.min(next.drawMode, next.stock.length);
        const drawn = next.stock.splice(next.stock.length - count, count);
        for (const card of drawn) {
            next.waste.push({ ...card, faceUp: true });
        }
        next.moves++;
        return next;
    }

    if (next.waste.length > 0) {
        next.stock = next.waste
            .slice()
            .reverse()
            .map((card) => ({ ...card, faceUp: false }));
        next.waste = [];
        next.moves++;
        return next;
    }

    return state;
}

export function canSelectSource(state: GameState, from: MoveSource): boolean {
    return getMoveStack(state, from) !== null;
}

export function getMovingCards(state: GameState, from: MoveSource): Card[] | null {
    return getMoveStack(state, from);
}

export function pileRefsEqual(left: PileRef, right: PileRef): boolean {
    return isSamePile(left, right);
}

export function applyMove(
    state: GameState,
    from: MoveSource,
    to: MoveDestination,
): GameState | null {
    if (isSamePile(from.pile, to.pile)) {
        return null;
    }

    const movingCards = getMoveStack(state, from);
    if (!movingCards) {
        return null;
    }

    const destinationPile = getPile(state, to.pile);

    if (to.pile.kind === 'tableau') {
        if (!canPlaceOnTableau(movingCards, destinationPile)) {
            return null;
        }
    } else if (to.pile.kind === 'foundation') {
        if (movingCards.length !== 1 || !canPlaceOnFoundation(movingCards[0]!, destinationPile)) {
            return null;
        }
    } else {
        return null;
    }

    const next = cloneGameState(state);
    const sourcePile = getPile(next, from.pile);
    const targetPile = getPile(next, to.pile);
    const movedCards = sourcePile.splice(from.cardIndex);
    targetPile.push(...movedCards);

    if (from.pile.kind === 'tableau') {
        revealTableauTop(sourcePile);
    }

    next.moves++;
    return next;
}

export function getValidDestinations(state: GameState, from: MoveSource): MoveDestination[] {
    const movingCards = getMoveStack(state, from);
    if (!movingCards) {
        return [];
    }

    const destinations: MoveDestination[] = [];

    for (let index = 0; index < TABLEAU_PILE_COUNT; index++) {
        const pile: PileRef = { kind: 'tableau', index };
        if (isSamePile(from.pile, pile)) {
            continue;
        }
        if (canPlaceOnTableau(movingCards, getPile(state, pile))) {
            destinations.push({ pile });
        }
    }

    if (movingCards.length === 1) {
        for (let index = 0; index < FOUNDATION_PILE_COUNT; index++) {
            const pile: PileRef = { kind: 'foundation', index };
            if (isSamePile(from.pile, pile)) {
                continue;
            }
            if (canPlaceOnFoundation(movingCards[0]!, getPile(state, pile))) {
                destinations.push({ pile });
            }
        }
    }

    return destinations;
}

export function isWon(state: GameState): boolean {
    return state.foundations.every((foundation) => foundation.length === 13);
}

export function cardLabel(rank: Rank): string {
    switch (rank) {
        case 1:
            return 'A';
        case 11:
            return 'J';
        case 12:
            return 'Q';
        case 13:
            return 'K';
        default:
            return String(rank);
    }
}

export function suitSymbol(suit: Suit): string {
    switch (suit) {
        case 'hearts':
            return '♥';
        case 'diamonds':
            return '♦';
        case 'clubs':
            return '♣';
        case 'spades':
            return '♠';
    }
}
