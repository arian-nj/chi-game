import { useBackendHealth } from '@/composables/use-backend-health';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { getMe } from '@/gen/account/v1/account-AccountService_connectquery';
import { computed } from 'vue';
import { useConnectQuery } from './use-connect-query';

export function useGuestProfile() {
  const { isBackendHealthy } = useBackendHealth();
  const { token } = useGuestAuth();

  const query = useConnectQuery(getMe, undefined, {
    enabled: computed(() => isBackendHealthy.value && Boolean(token.value)),
    staleTime: 5 * 60_000,
    retry: false,
  });

  const username = computed(() => query.data.value?.account?.username ?? null);

  const displayUsername = computed(() =>
    username.value ? `@${username.value}` : null,
  );

  return {
    username,
    displayUsername,
    isLoading: query.isLoading,
  };
}
