<script setup lang="ts">
import { AccountService } from '../gen/account/v1/account_pb';
import { createClient } from "@connectrpc/connect";
import { useQuery } from '@tanstack/vue-query';
import { authTransport } from '../lib/transport';

const { isPending, error, data } = useQuery({
  queryKey: ['me'],
  staleTime: 0,
  retryDelay: 10000,
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})

// Log error when it occurs
if (error.value) {
  console.error('Error fetching user data:', error.value)
}

</script>

<template>
  <div class="font-extrabold text-lg ">
    <span v-if="isPending">Loading...</span>
    <span v-else-if="error">Error</span>
    <h1 v-else>@{{ data?.account?.username }}</h1>
  </div>

</template>