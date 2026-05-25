<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import AddFriendButton from '../components/friends/AddFriendButton.vue';
import AddFriendsPage from '../components/friends/AddFriendsPage.vue';
import { useRoute, useRouter } from 'vue-router';
import { createClient } from '@connectrpc/connect';
import { ChatService } from '../gen/chat/v1/chat_pb';
import { authTransport } from '../lib/transport';
import { useQuery, useQueryClient } from '@tanstack/vue-query';

const scrollContainerRef = ref<HTMLDivElement | null>(null)

onMounted(() => {
	if (scrollContainerRef.value) {
		scrollContainerRef.value.addEventListener("scroll", handleScroll)
	}
})
onUnmounted(() => {
	if (scrollContainerRef.value) {
		scrollContainerRef.value.removeEventListener("scroll", handleScroll)
	}
})
const scrollTop = ref<number>(0)
function handleScroll(event: Event) {
	const target = event.target as HTMLDivElement
	scrollTop.value = target.scrollTop
}

// Add Friend
const router = useRouter();
const route = useRoute();


// Read from query parameter
const showAddFriendPage = computed({
	get: () => route.query.showAddFriend === 'true',
	set: (value: boolean) => {
		if (value) {
			router.replace({ query: { ...route.query, showAddFriend: 'true' } });
		} else {
			const { showAddFriend, ...restQuery } = route.query;
			router.replace({ query: restQuery });
		}
	}
});

function showFriendsPage() {
	showAddFriendPage.value = true
}
function hideFriendsPage() {
	showAddFriendPage.value = false
}

const queryClient = useQueryClient()
// query every 2 seconds
setInterval(() => {
	queryClient.invalidateQueries({ queryKey: ['chats'] })
}, 2_000)

const { data: chatsData } = useQuery({
	queryKey: ['chats'],
	queryFn: async () => {
		const client = createClient(ChatService, authTransport)
		const data = await client.getAllChats({})
		return data
	}
})

function goToChat(chatRoomId: bigint) {
	router.push(`/chat/${chatRoomId}`)
}
</script>

<template>
<div class="bg-custom-blue h-full">
		<div class="h-14 flex items-center gap-4 px-2">
			<h1 class="text-gray-200 font-extrabold text-4xl whitespace-nowrap shrink-0">Chi Game</h1>
		<div class="flex-1"></div>
		<AddFriendButton :scroll-top="scrollTop" class="shrink-0" @is-clicked="showFriendsPage">
		</AddFriendButton>
	</div>
	
	<div v-if="chatsData" class="bg-custom-lite-blue lg:max-w-1/2 mx-auto">
		<div v-for="chat in chatsData.chats" :key="chat.chatRoomId.toString()" @click="goToChat(chat.chatRoomId)" 
		class=" flex items-center py-2 pl-2 cursor-pointer">

			<!-- Avatar Placeholder -->
			<div class="w-12 h-12 rounded-full bg-linear-to-br from-blue-500 to-purple-500 flex items-center justify-center text-white font-semibold text-xl">
				{{ chat.otherPersonName.charAt(0).toUpperCase() }}
			</div>
			<div class="pl-4">
				<p class="text-gray-100 text-xl font-bold">{{ chat.otherPersonName }}</p>
				<p class="text-sm text-gray-300">last message here</p>
			</div>
		</div>
	</div>
	
	<Transition>
		<div v-if="showAddFriendPage">
			<AddFriendsPage @back-clicked="hideFriendsPage" />
		</div>
	</Transition>
	
</div>
</template>
