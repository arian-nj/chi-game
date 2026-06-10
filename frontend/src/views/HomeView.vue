<script setup lang="ts">
import GamesBoard from '@/components/GamesBoard.vue';
import HeaderComponent from '../components/header/HeaderComponent.vue';

</script>

<template>
  <div class="home-root lg:w-80/100 m-auto">
    <HeaderComponent />
    <GamesBoard />
  </div>
</template>

<style scoped>

.home-root {
  border: 1px solid red;
  min-height: 100vh;
  min-height: 100dvh;
}


</style>
