<script setup lang="ts">
import sudokuLogo from '@/assets/games/sudoku/sudoku_logo.svg';
</script>

<template>
    <section class="help-section" id="how-to-play-sudoku">
        <div class="help-window">
            <div class="help-titlebar">
                <div class="help-titlebar-left">
                    <img
                        :src="sudokuLogo"
                        alt=""
                        class="help-icon"
                        width="22"
                        height="22"
                    >
                    <span class="help-title">How to Play</span>
                </div>
            </div>
            <div class="help-body">
                <p>
                    Fill the <strong>9×9 grid</strong> so every row, column, and
                    <strong>3×3 box</strong> contains the digits <strong>1 through 9</strong>
                    without repeating. Some numbers are given at the start and cannot be changed.
                </p>
                <p>
                    Click a cell and type a number (<strong>1–9</strong>) to fill it in. Use
                    <strong>Backspace</strong> or <strong>Delete</strong> to clear your entry.
                    Conflicting numbers are highlighted in red.
                </p>
                <p>
                    Choose a difficulty from the Game menu — <strong>Easy</strong>,
                    <strong>Medium</strong>, or <strong>Hard</strong> — for more or fewer starting
                    clues. Start a <strong>New Game</strong> anytime to get a fresh puzzle.
                </p>
            </div>
        </div>
    </section>
</template>

<style scoped>
.help-section {
    display: flex;
    justify-content: center;
    margin-top: 1.5rem;
    padding: 0 0.75rem 2rem;
}

.help-window {
    width: min(100%, 36rem);
    border: 1px solid #0054e3;
    border-radius: 10px 10px 12px 12px;
    overflow: hidden;
    box-shadow:
        1px 1px 0 #000,
        0 8px 32px rgba(0, 20, 60, 0.35),
        0 2px 8px rgba(0, 0, 0, 0.2);
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
}

.help-titlebar {
    display: flex;
    align-items: center;
    padding: 6px 8px;
    background: linear-gradient(
        180deg,
        #0997ff 0%,
        #0053ee 8%,
        #0050ee 40%,
        #06f 88%,
        #06f 93%,
        #005eff 95%,
        #003ddb 96%,
        #003ddb 100%
    );
}

.help-titlebar-left {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
}

.help-icon {
    width: 22px;
    height: 22px;
    flex-shrink: 0;
    filter: drop-shadow(1px 1px 1px rgba(0, 0, 0, 0.45));
}

.help-title {
    font-size: 0.95rem;
    font-weight: bold;
    color: #fff;
    text-shadow: 1px 1px #000;
    line-height: 1;
}

.help-body {
    padding: 1.125rem 1.25rem 1.25rem;
    background: #ece9d8;
    border-top: 1px solid #0054e3;
    box-shadow: inset 1px 1px 0 #fff;
}

.help-body p {
    margin: 0;
    font-size: clamp(0.9rem, 2.5vw, 1rem);
    line-height: 1.6;
    font-weight: 400;
    color: #000;
}

.help-body p + p {
    margin-top: 0.75rem;
}

.help-body strong {
    font-weight: 700;
}
</style>
