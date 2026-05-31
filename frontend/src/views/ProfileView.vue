<!-- get and show others profile page --> 
<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import { AccountService } from '../gen/account/v1/account_pb';
import { authTransport } from '../lib/transport';
import { useQuery } from '@tanstack/vue-query';
import { ref } from 'vue';
import { useRoute } from 'vue-router';

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

</script>

<template>
    <div class="flex flex-col h-screen bg-custom-blue items-center">
        <h1 v-if="personIsPending">Loading...</h1>
        <h1 v-else-if="personError">Error</h1>

            <FriendStateButton :personId="personId" />
            <div class="bg-custom-lite-blue rounded-lg px-6 py-1 w-10/12 mt-4">
                <p class="text-xl font-bold ">@{{ personData?.account?.username }}</p>
                <p class="text-sm text-gray-400 font-medium">Username</p>
            </div>
    </div>
</template>
