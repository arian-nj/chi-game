<script setup lang="ts">
import { ref } from 'vue'

const count = ref(0)
</script>

<template>
  <section id="center">
    <button class="counter" @click="count++">Count is {{ count }}</button>
  </section>


</template>
