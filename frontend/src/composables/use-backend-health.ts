import { healthCheck } from '@/gen/healthcheck/v1/healthcheck-HealthcheckService_connectquery';
import { HealthType } from '@/gen/healthcheck/v1/healthcheck_pb';
import { computed } from 'vue';
import { useConnectQuery } from './use-connect-query';

export function useBackendHealth() {
  const query = useConnectQuery(healthCheck, undefined, {
    staleTime: 30_000,
    refetchInterval: 60_000,
    retry: 1,
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
