<script setup lang="ts">


const props = defineProps({
    scrollTop: {
        type : Number,
        required :true
    },
})

function onClick() {
    alert("clicked")
}

// const threshold = ref<number>(60)
// function updateThreshold() {
//     threshold.value = Math.max(60,window.innerHeight)
// }
// onMounted(()=>{
//     updateThreshold()
//     window.addEventListener('resize',updateThreshold)
// })
// onUnmounted(()=>{
//     window.removeEventListener('resize',updateThreshold)
// })

</script>

<template>
    <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 scale-75"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition-all duration-300 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-75"
        mode="out-in"
    >

        <button v-if="props.scrollTop < 60" @click="onClick" :class="[`
            w-12 h-12 rounded-full
            text-6xl flex items-center justify-center`,
        ]">
            <svg class="w-12 h-12 text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M12 4v16m8-8H4"></path>
            </svg>
        </button>
    </Transition>
</template>