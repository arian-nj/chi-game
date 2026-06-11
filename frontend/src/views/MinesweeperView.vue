<script setup lang="ts">
import MinesweeperGame from './MinesweeperGame.vue';


</script>

<template>
    <div style="" alt="windows xp bliss"
        class="min-h-screen w-full bg-cover bg-center bg-no-repeat"
    >
        <h1 class="text-center text-4xl font-bold pb-4">ChiGame.Site</h1>
        <div class="bg-[#0055e8] px-1 pb-1 m-auto w-fit rounded-sm">
            <!-- title name and logo -->
            <div class="bg-[#0055e8] py-1 flex">
                <img src="/public/games/minesweeper/minesweeper_logo.png" alt="windows xp bliss" class="w-6 h-6">
                <h2 class="font-bold xp-text-shadow">Minesweeper</h2>
            </div>
            <!-- game and help -->
            <div class="bg-[#e3e0d0] text-black text-sm w-full flex gap-4">
                <button class="hover:bg-[#2563eb] p-0.5">Game</button>
                <a href="#how-to-play-minesweeper" class="hover:bg-[#2563eb] p-0.5">Help</a>
            </div>
            <!-- minesweeper game -->
            <MinesweeperGame :cell-height="8" :cell-width="8" :mine-count="10" />
        </div>
    </div>
    
    <div class="bg-gray-400" id="how-to-play-minesweeper">
        <h1>How to play minesweeper</h1>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
        <p>Minesweeper is a game where you need to find all the mines on the board. You can do this by clicking on the squares. If you click on a mine, you lose. If you click on a square that is not a mine, you can see the number of mines around it. You can use this to help you find the mines.</p>
    
    </div>
    
</template>

<style scoped>
.xp-text-shadow {
    text-shadow: 1.5px 1.5px 1.5px #000;
}
</style>