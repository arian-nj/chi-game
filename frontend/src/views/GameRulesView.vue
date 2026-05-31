<script setup lang="ts">

</script>

<template>
  <div>
    <section class="max-w-2xl mx-auto p-6 bg-custom-lite-blue/40 rounded-2xl border border-white/10 shadow-md">
  <h2 class="text-3xl font-extrabold mb-6 tracking-wide text-center bg-gradient-to-r from-white to-blue-200 bg-clip-text text-transparent">
    How to Play XO (Tic Tac Toe)
  </h2>

  <div class="space-y-6 text-blue-100 leading-relaxed">
    <!-- Objective Section -->
    <div>
      <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
        🎯 Objective
      </h3>
      <p>
        The goal of the game is to be the first player to align <strong>three of your symbols</strong> in a continuous straight line. This line can be horizontal, vertical, or diagonal.
      </p>
    </div>

    <!-- Core Rules Section -->
    <div>
      <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
        📜 Gameplay Rules
      </h3>
      <ul class="list-disc pl-5 space-y-2">
        <li>The match is played on a classic <strong>3x3 grid</strong> consisting of 9 empty squares.</li>
        <li>Player 1 is assigned the <strong>X</strong> symbol, and Player 2 (or the AI) takes the <strong>O</strong> symbol.</li>
        <li>Players take turns placing their respective marks into any remaining unoccupied square.</li>
        <li>Once a mark is placed, it cannot be moved, overwritten, or removed from the board.</li>
      </ul>
    </div>

    <!-- Winning Section -->
    <div>
      <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
        🏆 Winning & Draws
      </h3>
      <ul class="list-disc pl-5 space-y-2">
        <li><strong>Victory:</strong> The instant a player successfully arranges 3 marks in a row, column, or diagonal, the game ends and that player wins.</li>
        <li><strong>Draw (Cats Game):</strong> If all 9 squares on the grid are completely filled and neither player has achieved a line of 3, the game ends in a stalemate/tie.</li>
      </ul>
    </div>

    <!-- Pro Strategy Tip -->
    <blockquote class="border-l-4 border-green-400 bg-green-500/10 p-4 rounded-r-xl mt-4">
      <span class="font-bold text-white block mb-1">💡 Quick Strategy Hint:</span>
      If you go first, controlling the center square or the outer corners gives you the highest mathematical probability to force a winning trap!
    </blockquote>
  </div>
</section>
  </div>
</template>