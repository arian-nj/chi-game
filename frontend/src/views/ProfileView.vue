<!-- get and show others profile page --> 
<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import { AccountService } from '../gen/account/v1/account_pb';
import { authTransport } from '../lib/transport';
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { FriendshipStatus, FriendsService } from '../gen/friends/v1/friends_pb';
import { useToast } from '../components/Toast.vue';

const routes = useRoute()
const pid = Number(routes.params.id)

const personId = ref(pid)

const { isPending: personIsPending, error: personError, data: personData } = useQuery({
    queryKey: [`person-${personId.value}`],
    queryFn: async () => {
        const client = createClient(AccountService, authTransport)
        const data = await client.getPerson({ id: BigInt(personId.value) })
        return data
    }
})

const fqueryKey = `person-friendship-${personId.value}`

const {isPending:fIsPending,error:fError,data:fData} = useQuery({
    queryKey: [fqueryKey],
    queryFn: async () => {
        const client = createClient(FriendsService,authTransport)
        const data = await client.getFriendshipStatus({ toPerson:BigInt(personId.value)})
        return data
    }
})

const clientQuery = useQueryClient()

function invalidateFriendshipQuery() {
    clientQuery.invalidateQueries({queryKey:[fqueryKey]})
}

const {toast} = useToast()

const sendFriendRequest = async () => {
    const client = createClient(FriendsService,authTransport)
    const data = await client.sendFriendReq({toPerson:BigInt(personId.value)})
    if (data.success) {
        toast.success("درخواست دوستی ارسال شد")
    }else {
        toast.warning("خطا دوباره امتحان کن")
    }
    invalidateFriendshipQuery()
}

const sendMessage = () => {
    alert('Opening message dialog...');
}

const cancelFriendRequest = async () => {
    const client = createClient(FriendsService,authTransport)
    const data = await client.cancelFriendReq({toPerson:BigInt(personId.value)})
    if (data.success) {
        toast.success("درخواست حذف شد")
    }else {
        toast.warning("خطا دوباره امتحان کن")
    }
    invalidateFriendshipQuery()
}

const acceptFriendRequest = async () => {
    const client = createClient(FriendsService,authTransport)
    const data = await client.acceptFriendReq({fromPerson:BigInt(personId.value)})
    if (data.success) {
        toast.success("درخواست قبول شد")
    }else {
        toast.warning("خطا دوباره امتحان کن")
    }
    invalidateFriendshipQuery()
}

</script>

<template>
    <div class="flex flex-col h-screen bg-custom-blue items-center">
        <h1 v-if="personIsPending">Loading...</h1>
        <h1 v-else-if="personError">Error</h1>

        <div class="flex flex-col items-center justify-center pt-10 max-w-lg w-full" v-else-if="personData?.account">
            <img src="../assets/profile.jpg" alt="profile_pic" class="rounded-full w-36">
            <h1 class="text-3xl font-bold">{{ personData?.account?.displayName }}</h1>

            <!-- Friendship Status Buttons -->
            <div class="w-1/2 mt-2">
                <!-- Loading State -->
                <button v-if="fIsPending" class="bg-blue-900/30 w-full h-16 rounded-xl text-2xl font-extrabold" disabled>
                    <p>Loading...</p>
                </button>

                <!-- Error State -->
                <button v-else-if="fError" class="bg-red-900/30 w-full h-16 rounded-xl text-2xl font-extrabold" disabled>
                    <p>Error</p>
                </button>

                <!-- Nothing State - Send Friend Request -->
                <button v-else-if="fData?.fstatus == FriendshipStatus.NOTHING" 
                        @click="sendFriendRequest"
                        class="bg-blue-600 hover:bg-blue-700 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
                    <p>دعوت به دوستی</p>
                </button>

                <!-- Friend State - Send Message -->
                <button v-else-if="fData?.fstatus == FriendshipStatus.FRIEND" 
                        @click="sendMessage"
                        class="bg-green-600 hover:bg-green-700 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
                    <p>پیام دادن</p>
                </button>

                <!-- Requested State - Cancel Request -->
                <button v-else-if="fData?.fstatus == FriendshipStatus.REQUESTED" 
                        @click="cancelFriendRequest"
                        class="bg-yellow-600 hover:bg-yellow-700 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
                    <p>پس گرفتن درخواست</p>
                </button>

                <!-- Received Request State - Accept Request -->
                <button v-else-if="fData?.fstatus == FriendshipStatus.RECEIVED_REQUEST" 
                        @click="acceptFriendRequest"
                        class="bg-purple-600 hover:bg-purple-700 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
                    <p>قبول درخواست</p>
                </button>
            </div>

            <div class="bg-custom-lite-blue rounded-lg px-6 py-1 w-10/12 mt-4">
                <p class="text-xl font-bold ">@{{ personData?.account?.username }}</p>
                <p class="text-sm text-gray-400 font-medium">Username</p>
            </div>
        </div>
    </div>
</template>
