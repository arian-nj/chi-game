<script setup lang="ts">
import { usePageSeo } from '@/composables/use-page-seo';

usePageSeo();
</script>

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
