import { isRed, type Card } from './types';

export function getMovableStack(pile: Card[], fromIndex: number): Card[] | null {
    if (fromIndex < 0 || fromIndex >= pile.length) {
        return null;
    }

    const stack = pile.slice(fromIndex);
    if (!stack.every((card) => card.faceUp)) {
        return null;
    }

    if (!isValidTableauSequence(stack)) {
        return null;
    }

    return stack;
}

export function isValidTableauSequence(cards: Card[]): boolean {
    if (cards.length === 0) {
        return false;
    }

    for (let i = 0; i < cards.length - 1; i++) {
        const below = cards[i]!;
        const above = cards[i + 1]!;
        if (below.rank !== above.rank + 1) {
            return false;
        }
        if (isRed(below.suit) === isRed(above.suit)) {
            return false;
        }
    }

    return true;
}

export function canPlaceOnTableau(cards: Card[], targetPile: Card[]): boolean {
    if (cards.length === 0) {
        return false;
    }

    if (!isValidTableauSequence(cards)) {
        return false;
    }

    const movingBottom = cards[0]!;

    if (targetPile.length === 0) {
        return movingBottom.rank === 13;
    }

    const targetTop = targetPile[targetPile.length - 1]!;
    return (
        movingBottom.rank === targetTop.rank - 1 &&
        isRed(movingBottom.suit) !== isRed(targetTop.suit)
    );
}

export function canPlaceOnFoundation(card: Card, foundationPile: Card[]): boolean {
    if (!card.faceUp) {
        return false;
    }

    if (foundationPile.length === 0) {
        return card.rank === 1;
    }

    const top = foundationPile[foundationPile.length - 1]!;
    return card.suit === top.suit && card.rank === top.rank + 1;
}
