<script setup lang="ts">
import { createClient } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { AccountService } from '../gen/account/v1/account_pb';
import { authTransport } from '../lib/transport';


const {isPending,error,data} = useQuery({
    queryKey:['me'],
    staleTime:0,
    retryDelay: 10000,
    queryFn: async () => {
        const client = createClient(AccountService, authTransport)
        const data = await client.getMe({})
        return data
    }
})
</script>

<template>
    <div class="flex flex-col items-center justify-center h-screen bg-custom-blue">
        <h1 v-if="isPending">Loading...</h1>
        <h1 v-else-if="error">Error</h1>

        <div v-else-if="data?.account">
            <h1 class="text-2xl font-bold">{{ data?.account?.username }}</h1>
            <h1 class="text-2xl font-bold">name {{ data?.account?.displayName }}</h1>
        </div>
        

        <!-- <img :src="data?.account?.profilePicture" alt="Profile Picture" class="w-24 h-24 rounded-full"> -->

    </div>
</template>