<script setup lang="ts">
import MinesweeperGame from './MinesweeperGame.vue';
</script>

<template>
    <div class="minesweeper-desktop min-h-screen w-full flex items-start justify-center pt-8">
        <div class="xp-window">
            <!-- title bar -->
            <div class="xp-titlebar">
                <div class="xp-titlebar-left">
                    <img
                        src="/games/minesweeper/minesweeper_logo.png"
                        alt=""
                        class="xp-titlebar-icon"
                    >
                    <span class="xp-titlebar-text">Minesweeper</span>
                </div>
                <div class="xp-titlebar-buttons">
                    <button type="button" class="xp-win-btn xp-win-btn-min" aria-label="Minimize" />
                    <button type="button" class="xp-win-btn xp-win-btn-max" aria-label="Maximize" />
                    <button type="button" class="xp-win-btn xp-win-btn-close" aria-label="Close" />
                </div>
            </div>

            <!-- menu bar -->
            <div class="xp-menubar">
                <button type="button" class="xp-menu-item">Game</button>
                <a href="#how-to-play-minesweeper" class="xp-menu-item">Help</a>
            </div>

            <!-- game -->
            <MinesweeperGame :cell-height="9" :cell-width="9" :mine-count="10" />
        </div>
    </div>

    <div class="bg-gray-400" id="how-to-play-minesweeper">
        <h1>How to play minesweeper</h1>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
    </div>
</template>

<style scoped>
.minesweeper-desktop {
    background: url('/games/minesweeper/bliss.jpg') center / cover no-repeat;
}

.xp-window {
    border: 1px solid #0054e3;
    border-radius: 8px 8px 0 0;
    box-shadow: 1px 1px 0 #000;
    overflow: hidden;
    font-family: Tahoma, 'MS Sans Serif', sans-serif;
}

.xp-titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 30px;
    padding: 0 3px 0 4px;
    background: linear-gradient(180deg, #0997ff 0%, #0053ee 8%, #0050ee 40%, #06f 88%, #06f 93%, #005eff 95%, #003ddb 96%, #003ddb 100%);
    border-radius: 8px 8px 0 0;
}

.xp-titlebar-left {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 0;
}

.xp-titlebar-icon {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
}

.xp-titlebar-text {
    color: #fff;
    font-size: 11px;
    font-weight: bold;
    text-shadow: 1px 1px #000;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.xp-titlebar-buttons {
    display: flex;
    gap: 2px;
    flex-shrink: 0;
}

.xp-win-btn {
    width: 21px;
    height: 21px;
    border: none;
    padding: 0;
    cursor: default;
    border-radius: 3px;
}

.xp-win-btn-min,
.xp-win-btn-max {
    background: linear-gradient(180deg, #3c8cfd 0%, #3b8bfc 50%, #3a8afb 51%, #3586f8 100%);
    border: 1px solid #003c74;
    box-shadow: inset 1px 1px 0 #6ba8ff;
    position: relative;
}

.xp-win-btn-min::after {
    content: '';
    position: absolute;
    left: 5px;
    right: 5px;
    top: 9px;
    height: 2px;
    background: #fff;
}

.xp-win-btn-max::after {
    content: '';
    position: absolute;
    left: 5px;
    top: 5px;
    width: 9px;
    height: 9px;
    border: 2px solid #fff;
    box-sizing: border-box;
}

.xp-win-btn-close {
    background: linear-gradient(180deg, #e88a8f 0%, #e4787d 50%, #e36b70 51%, #d94a4f 100%);
    border: 1px solid #8b0000;
    box-shadow: inset 1px 1px 0 #f0a8ab;
    position: relative;
}

.xp-win-btn-close::before,
.xp-win-btn-close::after {
    content: '';
    position: absolute;
    left: 4px;
    right: 4px;
    top: 9px;
    height: 2px;
    background: #fff;
}

.xp-win-btn-close::before {
    transform: rotate(45deg);
}

.xp-win-btn-close::after {
    transform: rotate(-45deg);
}

.xp-menubar {
    display: flex;
    gap: 0;
    background: #ece9d8;
    padding: 2px 0;
    border-bottom: 1px solid #aca899;
}

.xp-menu-item {
    background: none;
    border: none;
    color: #000;
    font-size: 11px;
    font-family: inherit;
    padding: 1px 6px;
    cursor: default;
    text-decoration: none;
}

.xp-menu-item:hover {
    background: #316ac5;
    color: #fff;
}
</style>
