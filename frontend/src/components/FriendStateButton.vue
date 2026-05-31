<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import { authTransport } from '../lib/transport';
import { useToast } from '../components/Toast.vue';
import { useQuery, useQueryClient } from '@tanstack/vue-query';
import { ref } from 'vue';
import { AccountService } from '../gen/account/v1/account_pb';
import { FriendshipStatus, FriendsService } from '../gen/friends/v1/friends_pb';
import { useRouter } from 'vue-router';


const personId = ref(0)

const { isPending: personIsPending, error: personError, data: personData } = useQuery({
    queryKey: [`person-${personId.value}`],
    queryFn: async () => {
        const client = createClient(AccountService, authTransport)
        const data = await client.getPerson({ id: BigInt(personId.value) })
        return data
    }
})

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
const router = useRouter()

const sendMessage = () => {
    router.push(`/chat/${personId.value}`)
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
const clientQuery = useQueryClient()

function invalidateFriendshipQuery() {
    clientQuery.invalidateQueries({queryKey:[fqueryKey]})
}


const fqueryKey = `person-friendship-${personId.value}`

const {isPending:fIsPending,error:fError,data:fData} = useQuery({
    queryKey: [fqueryKey],
    queryFn: async () => {
        const client = createClient(FriendsService,authTransport)
        const data = await client.getFriendshipStatus({ toPerson:BigInt(personId.value)})
        return data
    }
})


</script>

<template>
<div class="flex flex-col items-center justify-center pt-10 max-w-lg w-full" v-if="personData?.account">
    <img src="../assets/profile.jpg" alt="profile_pic" class="rounded-full w-36">
    <h1 class="text-3xl font-bold">{{ personData?.account?.displayName }}</h1>

    <!-- Friendship Status Buttons -->
    <div class="w-1/2 mt-2">
        <!-- Loading State -->
        <button v-if="fIsPending" class="bg-gray-500/60 w-full h-16 rounded-xl text-2xl font-extrabold" disabled>
            <p>Loading...</p>
        </button>

        <!-- Error State -->
        <button v-else-if="fError" class="bg-red-700/50 w-full h-16 rounded-xl text-2xl font-extrabold" disabled>
            <p>Error</p>
        </button>

        <!-- Nothing State - Send Friend Request -->
        <button v-else-if="fData?.fstatus == FriendshipStatus.NOTHING" 
                @click="sendFriendRequest"
                class="bg-cyan-500 hover:bg-cyan-600 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
            <p>دعوت به دوستی</p>
        </button>

        <!-- Friend State - Send Message -->
        <button v-else-if="fData?.fstatus == FriendshipStatus.FRIEND" 
                @click="sendMessage"
                class="bg-emerald-600 hover:bg-emerald-700 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
            <p>پیام دادن</p>
        </button>

        <!-- Requested State - Cancel Request -->
        <button v-else-if="fData?.fstatus == FriendshipStatus.REQUESTED" 
                @click="cancelFriendRequest"
                class="bg-orange-400 hover:bg-orange-500 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
            <p>پس گرفتن درخواست</p>
        </button>

        <!-- Received Request State - Accept Request -->
        <button v-else-if="fData?.fstatus == FriendshipStatus.RECEIVED_REQUEST" 
                @click="acceptFriendRequest"
                class="bg-fuchsia-500 hover:bg-fuchsia-600 w-full h-16 rounded-xl text-2xl font-extrabold transition-colors">
            <p>قبول درخواست</p>
        </button>
    </div>
</div>

</template>