<script setup lang="ts" >
import { createClient } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { HealthcheckService, HealthType } from '../gen/healthcheck/v1/healthcheck_pb';
import { rawTransport } from '../lib/transport';

const {isPending, error, data} = useQuery({
    queryKey : ["healthcheck"],
    staleTime: 0,
    queryFn: async () => {
        const client = createClient(HealthcheckService, rawTransport);
        const data = await client.healthCheck({});
        return data;
    }
});
</script>

<template>
    <div class="status-dot">
        <div 
            class="dot"
            :class="{
                'dot-green': data?.healthType == HealthType.OK,
                'dot-red': data?.healthType != HealthType.OK && !isPending && !error,
                'dot-yellow': isPending
            }"
        ></div>
        <span v-if="error" class="error-text">{{ error.message }}</span>
    </div>
</template>

<style scoped>
.status-dot {
    display: inline-flex;
    align-items: center;
    gap: 8px;
}

.dot {
    width: 14px;
    height: 14px;
    border-radius: 50%;
    display: inline-block;
}

.dot-green {
    background-color: #22c55e;
    box-shadow: 0 0 8px rgba(34, 197, 94, 0.5);
}

.dot-red {
    background-color: #ef4444;
    box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.dot-yellow {
    background-color: #eab308;
    box-shadow: 0 0 8px rgba(234, 179, 8, 0.5);
}

.error-text {
    color: #ef4444;
    font-size: 0.875rem;
}
</style>