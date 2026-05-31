<script setup lang="ts">
import { onMounted, ref } from 'vue';
import gsap from 'gsap'

const props = defineProps({
  gameName: String,
  selected: Boolean,
})

const emit = defineEmits(['choosed'])

const btnRef = ref<HTMLDivElement | null>()

onMounted(() => {
  if (btnRef.value) {
    gsap.from(btnRef.value, {
      scale: 0,
      opacity: 0.4,
      duration: .7,
    })
  }
})

</script>
<template>
  <button ref="btnRef" type="button" @click="emit('choosed')" :class="[
    `px-10 py-6 w-50 h-30 rounded-2xl text-2xl font-extrabold tracking-wide
         transition-colors duration-300 ease-in-out
         backdrop-blur-md`,
    props.selected
      ? 'bg-linear-to-br from-red-100/10 to-orange-100/10 text-white shadow-lg border-4 border-gray-300'
      : 'border-2 bg-gray-800/40 border-gray-600 text-gray-300 '
  ]">
    {{ props.gameName }}
  </button>
</template>
