export type Suit = 'hearts' | 'diamonds' | 'clubs' | 'spades';
export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13; // Ace=1, King=13
export type DrawMode = 1 | 3;
export type PileKind = 'stock' | 'waste' | 'foundation' | 'tableau';

export const SUITS: Suit[] = ['hearts', 'diamonds', 'clubs', 'spades'];
export const RANKS: Rank[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13];
export const TABLEAU_PILE_COUNT = 7;
export const FOUNDATION_PILE_COUNT = 4;

export interface Card {
    suit: Suit;
    rank: Rank;
    faceUp: boolean;
}

export interface GameState {
    stock: Card[];
    waste: Card[];
    foundations: Card[][];
    tableau: Card[][];
    drawMode: DrawMode;
    moves: number;
}

export type PileRef =
    | { kind: 'waste' }
    | { kind: 'foundation'; index: number }
    | { kind: 'tableau'; index: number };

export interface MoveSource {
    pile: PileRef;
    cardIndex: number;
}

export interface MoveDestination {
    pile: PileRef;
}

export function isRed(suit: Suit): boolean {
    return suit === 'hearts' || suit === 'diamonds';
}
