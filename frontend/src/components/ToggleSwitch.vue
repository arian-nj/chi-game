<script setup lang="ts">
import { computed, ref } from 'vue';

const props = defineProps<{
  optionOne: string;
  optionTwo: string;
}>();

const activeIndex = ref<1 | 2>(1);

const isFirstActive = computed(() => activeIndex.value === 1);

function selectOption(index: 1 | 2) {
  activeIndex.value = index;
}

const activeOptionNumber = computed(() => activeIndex.value);

defineExpose({
  optionNumber: activeOptionNumber,
});
</script>

<template>
  <div
    dir="ltr"
    class="relative inline-flex w-full min-w-56 rounded-2xl border border-white/10 bg-custom-deep-blue/90 p-1 shadow-[inset_0_2px_8px_rgba(0,0,0,0.35)]"
    role="group"
    :aria-label="`${props.optionOne} or ${props.optionTwo}`"
  >
    <span
      class="pointer-events-none absolute top-1 bottom-1 w-[calc(50%-0.25rem)] rounded-xl bg-white/90 shadow-md ring-1 ring-white/30 transition-[left,box-shadow] duration-300 ease-[cubic-bezier(0.34,1.2,0.64,1)]"
      :class="isFirstActive ? 'left-1' : 'left-[calc(50%+0.125rem)]'"
      aria-hidden="true"
    ></span>

    <button
      type="button"
      dir="auto"
      class="relative z-10 flex flex-1 cursor-pointer items-center justify-center rounded-xl px-4 py-2.5 text-lg font-bold tracking-wide transition-colors duration-200 select-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 focus-visible:ring-offset-2 focus-visible:ring-offset-custom-deep-blue"
      :class="isFirstActive ? 'text-custom-blue' : 'text-blue-100/80 hover:text-blue-50'"
      :aria-pressed="isFirstActive"
      @click="selectOption(1)"
    >
      {{ props.optionOne }}
    </button>

    <button
      type="button"
      dir="auto"
      class="relative z-10 flex flex-1 cursor-pointer items-center justify-center rounded-xl px-4 py-2.5 text-lg font-bold tracking-wide transition-colors duration-200 select-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/50 focus-visible:ring-offset-2 focus-visible:ring-offset-custom-deep-blue"
      :class="!isFirstActive ? 'text-custom-blue' : 'text-blue-100/80 hover:text-blue-50'"
      :aria-pressed="!isFirstActive"
      @click="selectOption(2)"
    >
      {{ props.optionTwo }}
    </button>
  </div>
</template>
