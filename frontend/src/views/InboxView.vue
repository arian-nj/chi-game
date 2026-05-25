<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import AddFriendButton from '../components/friends/AddFriendButton.vue';
import AddFriendsPage from '../components/friends/AddFriendsPage.vue';
import { useRoute, useRouter } from 'vue-router';
import { createClient } from '@connectrpc/connect';
import { ChatService } from '../gen/chat/v1/chat_pb';
import { authTransport } from '../lib/transport';
import { useQuery } from '@tanstack/vue-query';

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

const { data: chatsData } = useQuery({
	queryKey: ['chats'],
	queryFn: async () => {
		const client = createClient(ChatService, authTransport)
		const data = await client.getAllChats({})
		return data
	}
})

</script>

<template>
<div class="bg-custom-blue h-full">
		<div class="h-14 flex items-center gap-4 px-2">
			<h1 class="text-gray-200 font-extrabold text-4xl whitespace-nowrap shrink-0">Chi Game</h1>
		<div class="flex-1"></div>
		<AddFriendButton :scroll-top="scrollTop" class="shrink-0" @is-clicked="showFriendsPage">
		</AddFriendButton>
	</div>
	
	<div v-if="chatsData" class="bg-custom-lite-blue">
		<div v-for="chat in chatsData.chats" :key="chat.id.toString()">
			{{ chat.id }}
		</div>
	</div>
	
	<Transition>
		<div v-if="showAddFriendPage">
			<AddFriendsPage @back-clicked="hideFriendsPage" />
		</div>
	</Transition>
	
</div>
</template>
