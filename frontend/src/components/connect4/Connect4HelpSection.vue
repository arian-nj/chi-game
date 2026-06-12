<script setup lang="ts">
import connect4Logo from '@/assets/games/connect4/connect4_logo.svg';
</script>

<template>
    <section class="help-section" id="how-to-play-connect4">
        <div class="help-window">
            <div class="help-titlebar">
                <div class="help-titlebar-left">
                    <img
                        :src="connect4Logo"
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
                    Players take turns dropping colored discs into a <strong>7-column, 6-row</strong> grid.
                    Discs fall to the lowest open slot in the column you choose. The first player to connect
                    <strong>four discs in a row</strong> — horizontally, vertically, or diagonally — wins.
                    If the board fills with no winner, the game is a draw.
                </p>
                <p>
                    Click a column at the top of the board to drop your disc. You can only play in columns
                    that still have space.
                </p>
                <p>
                    Choose <strong>2 Players</strong> to play on the same device, or
                    <strong>vs Computer</strong> to play as <strong>red</strong> against the bot
                    (<strong>yellow</strong>). Difficulty can be set to Easy, Medium, or Hard from the
                    Game menu.
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

@media (max-width: 640px) {
    .help-section {
        margin-top: 0.75rem;
        padding: 0 0.25rem 1.25rem;
    }

    .help-window {
        width: 100%;
    }
}
</style>
