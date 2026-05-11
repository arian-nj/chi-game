<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import AddFriendButton from '../components/friends/AddFriendButton.vue';
import AddFriendsPage from '../components/friends/AddFriendsPage.vue';

const scrollContainerRef = ref<HTMLDivElement|null>(null)

onMounted(()=> {
    if (scrollContainerRef.value) {
        scrollContainerRef.value.addEventListener("scroll",handleScroll)
    }
})
onUnmounted(()=> {
    if (scrollContainerRef.value){
        scrollContainerRef.value.removeEventListener("scroll",handleScroll)
    }
})
const scrollTop = ref<number>(0)
function handleScroll(event:Event) {
    const target = event.target as HTMLDivElement
    scrollTop.value = target.scrollTop
}

// Add Friend
const showAddFriendPage = ref<boolean>(false)

function showFriendsPage() {
    showAddFriendPage.value = true
}
function hideFriendsPage() {
    showAddFriendPage.value = false
}

</script>

<template>
<div class="h-14 bg-custom-blue flex items-center gap-4 px-2">
    <h1 class="text-gray-200 font-extrabold text-4xl whitespace-nowrap shrink-0">Chi Game</h1>
    <div class="flex-1"></div>
    <AddFriendButton :scroll-top="scrollTop" class="shrink-0" @is-clicked="showFriendsPage"></AddFriendButton>
    <!-- <div class="flex-1 bg-custom-lite-blue/50 rounded-full flex items-center h-9/12 shadow-sm min-w-0">
        <button 
          type="submit" 
          class="pl-3 py-2 text-gray-400 transition-colors shrink-0 text-3xl"
          aria-label="Search"
        >
            <svg 
            class="w-6 h-6" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24" 
            xmlns="http://www.w3.org/2000/svg"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="4" 
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            ></path>
          </svg>
        </button>
      <input 
        type="text" 
        placeholder="Search Chats" 
        class="flex-1 pl-1.5 py-2 outline-none rounded-full text-gray-100/80 placeholder-gray-400/80 text-xl font-semibold"
      >
    </div> -->
</div>


<div ref="scrollContainerRef" class="h-screen relative overflow-y-auto bg-custom-blue">
    <div v-for="n in 50" :key="n" class="p-4">
        Friend {{ n }}
    </div>
</div>

<Transition>
    <div v-if="showAddFriendPage">
        <AddFriendsPage @back-clicked="hideFriendsPage"/>
    </div>
</Transition>

</template>