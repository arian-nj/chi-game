<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import AddFriendButton from '../components/friends/AddFriendButton.vue';
import AddFriendsPage from '../components/friends/AddFriendsPage.vue';
import { useRoute, useRouter } from 'vue-router';

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

</script>

<template>
	<div class="h-14 bg-custom-blue flex items-center gap-4 px-2">
		<h1 class="text-gray-200 font-extrabold text-4xl whitespace-nowrap shrink-0">Chi Game</h1>
		<div class="flex-1"></div>
		<AddFriendButton :scroll-top="scrollTop" class="shrink-0" @is-clicked="showFriendsPage">
		</AddFriendButton>
	</div>


	<div ref="scrollContainerRef" class="h-screen relative overflow-y-auto bg-custom-blue">
		<div v-for="n in 50" :key="n" class="p-4">
			Friend {{ n }}
		</div>
	</div>

	<Transition>
		<div v-if="showAddFriendPage">
			<AddFriendsPage @back-clicked="hideFriendsPage" />
		</div>
	</Transition>

</template>
