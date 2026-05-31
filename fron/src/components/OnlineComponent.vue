<script setup lang="ts" >
import { createClient } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import { HealthcheckService, HealthType } from '../gen/healthcheck/v1/healthcheck_pb';
import { rawTransport } from '../lib/transport';
import { watchEffect, computed } from 'vue';
import { useToast } from './Toast.vue';

const { toast } = useToast()

const { isPending, error, data, isFetching, isError } = useQuery({
    queryKey: ["healthcheck"],
    staleTime: Infinity,
    refetchInterval: 5000,
    retry: 3,
    retryDelay: 10000,
    queryFn: async () => {
        const client = createClient(HealthcheckService, rawTransport);
        const data = await client.healthCheck({});
        return data;
    }
});

// Compute the dot status more reliably
const dotStatus = computed(() => {
    // Still loading initial data
    if (isPending.value) {
        return 'yellow';
    }
    
    // Currently refetching but has existing data
    if (isFetching.value && data.value) {
        return 'green'; // Or 'yellow' - depending on your preference
    }
    
    // Error state
    if (isError.value || error.value) {
        return 'red';
    }
    
    // Has data
    if (data.value) {
        return data.value.healthType === HealthType.OK ? 'green' : 'red';
    }
    
    // Default fallback (should not happen often)
    return 'yellow';
});

// Watch for errors with debounce to avoid multiple toasts
let errorToastShown = false;
watchEffect(() => {
    if (error.value && !errorToastShown) {
        errorToastShown = true;
        toast.error("Health check failed");
        console.error("Health check failed:", {
            message: error.value.message,
            name: error.value.name,
            stack: error.value.stack
        });
        
        // Reset after a short delay to allow new errors to trigger toast
        setTimeout(() => {
            errorToastShown = false;
        }, 1000);
    }
});

// Reset error toast flag when error clears
watchEffect(() => {
    if (!error.value) {
        errorToastShown = false;
    }
});
</script>

<template>
    <div class="status-dot">
        <div 
            class="dot"
            :class="{
                'dot-green': dotStatus === 'green',
                'dot-red': dotStatus === 'red',
                'dot-yellow': dotStatus === 'yellow'
            }"
        ></div>

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
    transition: all 0.2s ease-in-out;
    min-width: 14px; /* Prevent shrinking */
    min-height: 14px; /* Prevent shrinking */
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

.status-text {
    font-size: 0.875rem;
    color: #6b7280;
}

.error-text {
    color: #ef4444;
}
</style>