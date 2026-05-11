<script setup lang="ts">
import BackButton from '../BackButton.vue';
import { createClient } from '@connectrpc/connect';
import { FriendsService, type SearchForUsernameResponse } from '../../gen/friends/v1/friends_pb';
import { authTransport } from '../../lib/transport';
import { ref, watch } from 'vue';


const emit = defineEmits<{
    backClicked:[]
}>()

const backClick = ()=>{emit("backClicked")}

const searchQuery = ref<string>('')
const searchResults = ref<SearchForUsernameResponse | null>(null)
const isSearching = ref(false)
const searchError = ref(null)

let activeAbortController: AbortController | null = null

async function performSearch(lookFor:string) {
    if (!lookFor.trim()) {
        searchResults.value = null
        return
    }

    if (activeAbortController) {
        activeAbortController.abort()
    }
    activeAbortController = new AbortController()
    isSearching.value = true

    try {
        const client = createClient(FriendsService,authTransport)
        const data = await client.searchForUsername({lookFor:lookFor},{
            signal: activeAbortController.signal
        })
        searchResults.value = data
    } catch (err:any) {
        if (err.name !== 'AbortError') {
        searchError.value = err
        console.error('Search failed', err)
        }
    } finally {
        isSearching.value = false
        activeAbortController = null
    }
}

let timeoutId: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, (newQuery) => {
  if (timeoutId) clearTimeout(timeoutId)
  timeoutId = setTimeout(() => performSearch(newQuery), 300)
})

</script>

<template>
<div class="bg-custom-blue w-screen h-screen z-50 fixed top-0 left-0 overflow-y-auto">
    <!-- Header with Back Button -->
    <div class="sticky top-0 bg-custom-blue/95 backdrop-blur-sm z-10">
        <div class="h-16 flex items-center">
            <BackButton @back-clicked="backClick" class="pl-4" />
        </div>
    </div>

    <!-- Search Box -->
    <div class="flex justify-center px-4 mt-4">
        <div class="flex-1 bg-custom-lite-blue rounded-full flex items-center max-w-2xl shadow-lg hover:shadow-xl transition-shadow duration-200">
            <button 
                type="submit" 
                class="pl-4 py-3 text-gray-400 transition-colors shrink-0 hover:text-gray-300"
                aria-label="Search"
            >
                <svg 
                    class="w-5 h-5" 
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24" 
                    xmlns="http://www.w3.org/2000/svg"
                >
                    <path 
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                    ></path>
                </svg>
            </button>

            <input 
                type="text" 
                placeholder="Search Chats..." 
                class="flex-1 pl-2 pr-4 py-3 outline-none rounded-full text-gray-100 placeholder-gray-400 text-base font-medium bg-transparent"
                v-model="searchQuery"
            >
        </div>
    </div>

    <!-- Content -->
    <div class="px-4 mt-6">
        <!-- Found Users Section -->
        <div v-if="searchResults?.FoundUsers?.length" class="max-w-2xl mx-auto">
            <h2 class="text-gray-400 text-sm font-semibold uppercase tracking-wider mb-3 px-2">
                Found Users
            </h2>
            
            <div class="space-y-2">
                <div 
                    v-for="person in searchResults.FoundUsers" 
                    class="group bg-custom-lite-blue/50 rounded-xl p-3 hover:bg-custom-lite-blue transition-all duration-200 cursor-pointer"
                >
                    <div class="flex items-center space-x-3">
                        <!-- Avatar Placeholder -->
                        <div class="w-10 h-10 rounded-full bg-linear-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-semibold text-sm">
                            {{ person.username.charAt(0).toUpperCase() }}
                        </div>
                        
                        <!-- User Info -->
                        <div class="flex-1 min-w-0">
                            <h3 class="text-gray-100 font-semibold text-base truncate">
                                {{ person.diplayName || `@${person.username}` }}
                            </h3>
                            <p class="text-gray-400 text-sm truncate">
                                @{{ person.username }}
                            </p>
                        </div>
                        
                        <!-- Action Button -->
                        <button class="px-4 py-1.5 bg-blue-500 hover:bg-blue-600 rounded-full text-white text-sm font-medium transition-colors duration-200">
                            Message
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- No Results State -->
        <div v-else-if="searchQuery && !searchResults?.FoundUsers?.length" class="text-center py-12">
            <div class="w-16 h-16 mx-auto mb-4 bg-custom-lite-blue rounded-full flex items-center justify-center">
                <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
            </div>
            <h3 class="text-gray-300 font-medium text-lg">No users found</h3>
            <p class="text-gray-400 text-sm mt-1">Try searching with a different username</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="!searchQuery" class="text-center py-12">
            <div class="w-16 h-16 mx-auto mb-4 bg-custom-lite-blue rounded-full flex items-center justify-center">
                <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
            </div>
            <h3 class="text-gray-300 font-medium text-lg">Search for users</h3>
            <p class="text-gray-400 text-sm mt-1">Find friends and start chatting</p>
        </div>
    </div>
</div>
</template>

<style scoped>
/* Smooth scroll behavior */
::-webkit-scrollbar {
    width: 6px;
}

::-webkit-scrollbar-track {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 10px;
}

::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.3);
    border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
    background: rgba(255, 255, 255, 0.5);
}
</style>