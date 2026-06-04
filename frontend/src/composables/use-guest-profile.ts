import { useBackendHealth } from '@/composables/use-backend-health';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { AccountService } from '@/gen/account/v1/account_pb';
import { createApiClient } from '@/libs/api-client';
import { useQuery } from '@tanstack/vue-query';
import { computed } from 'vue';

export function useGuestProfile() {
  const { isBackendHealthy } = useBackendHealth();
  const { token } = useGuestAuth();

  const client = createApiClient(AccountService);
  
  const query = useQuery({
    queryKey: ['account', 'me'],
    queryFn: ({ signal }) => client.getMe({}, { signal }),
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
