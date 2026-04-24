<script setup lang="ts" >
import { createClient } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { HealthcheckService, HealthType } from '../gen/healthcheck/v1/healthcheck_pb';
import { rawTransport } from '../lib/transport';


const {isPending,error,data} = useQuery({
    queryKey : ["healthcheck"],
    staleTime: 0,
    queryFn: async () => {
        const client = createClient(HealthcheckService,rawTransport)
        const data = await client.healthCheck({})
        return data
    }

})

</script>

<template>
    <span v-if="isPending">checking...</span>
    <span v-else-if="error">{{ error.message }}</span>
    <h1 v-else>{{ data?.healthType == HealthType.OK ? "online" : "offline"}}</h1>

</template>