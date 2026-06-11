<script setup lang="ts"></script>

<template>
  <div class="bliss-desktop">
    <RouterView />
  </div>
</template>

<style scoped>
.bliss-desktop {
    background: url('/public/bliss.jpg') center / cover no-repeat;
    box-sizing: border-box;
    padding: 2rem 0.5rem;
    overflow-x: hidden;
}

</style>
