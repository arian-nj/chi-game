<script setup lang="ts"></script>

<template>
  <div class="bliss-desktop">
    <RouterView />
  </div>
</template>

<style scoped>
.bliss-desktop {
    width: 100%;
    min-height: 100dvh;
    background: url('/public/bliss.jpg') center / cover no-repeat;
    box-sizing: border-box;
    overflow-x: hidden;
}

</style>
