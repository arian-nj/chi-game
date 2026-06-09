import { HealthcheckService, HealthType } from '@/gen/healthcheck/v1/healthcheck_pb';
import { createApiClient } from '@/libs/api-client';
import { useQuery } from '@tanstack/vue-query';
import { computed } from 'vue';

export function useBackendHealth() {
  const client = createApiClient(HealthcheckService);
  
  const query = useQuery({
    queryKey: ['backend-health'],
    queryFn: ({ signal }) => client.healthCheck({}, { signal }),
    staleTime: 60_000,
  });

  const isBackendHealthy = computed(
    () => query.data.value?.healthType === HealthType.OK,
  );

  const backendError = computed(() => {
    const err = query.error.value;
    return err instanceof Error ? err.message : err ? String(err) : null;
  });

  return {
    ...query,
    isBackendHealthy,
    backendError,
    isLoading: query.isLoading,
  };
}
