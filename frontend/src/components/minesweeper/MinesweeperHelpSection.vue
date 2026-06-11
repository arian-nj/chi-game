<script setup lang="ts">
import minesweeperLogo from '@/assets/games/minesweeper/minesweeper_logo.png';
</script>

<template>
    <section class="help-section" id="how-to-play-minesweeper">
        <div class="help-window">
            <div class="help-titlebar">
                <div class="help-titlebar-left">
                    <img
                        :src="minesweeperLogo"
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
                    Clear the board without hitting a mine. <strong>Left-click</strong> a square to
                    reveal it — if it's a mine, you lose. Safe squares show a number telling you
                    how many mines touch that square (including diagonals). Use those clues to find
                    where the mines are hidden.
                </p>
                <p>
                    <strong>Right-click</strong> a square to place or remove a flag on a suspected
                    mine. The counter at the top tracks how many mines are left to find. Reveal every
                    safe square to win.
                </p>
                <p>
                    Pick a difficulty from the Game menu — <strong>Beginner</strong> (9×9),
                    <strong>Intermediate</strong> (16×16), or <strong>Expert</strong> (30×16) — for
                    larger boards and more mines. Click the smiley face or choose
                    <strong>New Game</strong> to start over.
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
