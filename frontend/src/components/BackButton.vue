<script setup lang="ts">

const emit = defineEmits<{
    backClicked:[]
}>()

const backClick = ()=>{emit("backClicked")}

</script>

<template>
    <button
    type="button"
    aria-label="Go back"
    class="text-gray-200
     hover:text-gray-100
     transition-colors duration-200"
    @click="backClick"
  >
    <!-- Left Arrow SVG -->
    <svg 
      xmlns="http://www.w3.org/2000/svg" 
      class="w-8 h-8" 
      fill="none" 
      viewBox="0 0 24 24" 
      stroke="currentColor" 
      stroke-width="6"
    >
      <path 
        stroke-linecap="round" 
        stroke-linejoin="round" 
        d="M10 19l-7-7m0 0l7-7m-7 7h18" 
      />
    </svg>
  </button>

</template>