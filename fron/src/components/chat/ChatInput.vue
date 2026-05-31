<script setup lang="ts">
import { ref } from 'vue';

const text = ref("")

const emit = defineEmits<{
  submit: [text: string],
  inputClick: []
}>()

function onSubmit() {
  emit('submit', text.value)
  text.value = ""
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === "Enter" && !e.shiftKey) {
    e.preventDefault()
    onSubmit()
  }
}
</script>

<template>
  <div class="w-full">
    <form @submit.prevent="onSubmit">
      <label for="chat" class="sr-only">Your message</label>
      <div class="flex items-center px-3 py-2 rounded-t-lg">
        <textarea
          maxlength="256"
          id="chat"
          rows="1"
          v-model="text"
          dir="auto"
          @focus="$emit('inputClick')"
          @keydown="onKeydown"
          class="block mx-4 p-2.5 w-full text-gray-50 rounded-sm border border-gray-300 focus:ring-blue-500 focus:border-blue-500 bg-gray-800 text-xl"
          placeholder="Your message..."
        >
      </textarea>

        <button type="submit"
          class="inline-flex justify-center p-2 text-blue-600 rounded-full cursor-pointer hover:bg-blue-100">
          <svg class="w-7 h-7 rotate-90" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 18 20">
            <path
              d="m17.914 18.594-8-18a1 1 0 0 0-1.828 0l-8 18a1 1 0 0 0 1.157 1.376L8 18.281V9a1 1 0 0 1 2 0v9.281l6.758 1.689a1 1 0 0 0 1.156-1.376Z" />
          </svg>
          <span class="sr-only">Send message</span>
        </button>
      </div>
    </form>
  </div>
</template>
