<script setup lang="ts">
import AboutUsSection from '@/components/AboutUsSection.vue';
import GamesBoard from '@/components/GamesBoard.vue';
import HeaderComponent from '../components/header/HeaderComponent.vue';

</script>

<template>
  <div class="home-root lg:w-80/100 m-auto">
    <HeaderComponent />
    <GamesBoard />
    <AboutUsSection />
  </div>
</template>

<style scoped>

.home-root {
  /* border: 1px solid; */
  min-height: 100vh;
  min-height: 100dvh;
}


</style>
