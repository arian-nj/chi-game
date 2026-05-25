<script setup lang="ts">

// import profilePic from '@/assets/profile.jpg'

const props = defineProps({
  isMe: { required: true, type: Boolean },
  text: { required: true, type: String }
})

</script>

<template>

  <div :class="['flex items-end gap-2.5', props.isMe ? 'flex-row-reverse' : '']">
    <!-- <img :src="profilePic" class='w-10 h-10 rounded-full' alt="profile image" /> -->
    <div :class="['max-w-1/2 rounded-t-xl p-3 wrap-break-word', props.isMe ? 'bg-indigo-800 rounded-l-xl'
      : 'bg-gray-800 rounded-r-xl']">
      <p dir="auto">{{ props.text }}</p>
    </div>
  </div>
</template>
