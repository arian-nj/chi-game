import { useBackendHealth } from '@/composables/use-backend-health';
import { AuthService } from '@/gen/auth/v1/auth_pb';
import { createApiClient } from '@/libs/api-client';
import {
  readGuestDeviceId,
  readGuestToken,
  saveGuestSession,
} from '@/libs/guest-auth-storage';
import { computed, ref, watch } from 'vue';

// Shared across all useGuestAuth() callers (App + views).
const deviceId = ref(readGuestDeviceId());
const token = ref(readGuestToken());
const isGuest = computed(() => Boolean(token.value));

let ensureInFlight: Promise<void> | null = null;
let sessionWatchStarted = false;

async function ensureGuestSession(isBackendHealthy: boolean) {
  if (!isBackendHealthy) {
    return;
  }

  try {
    const client = createApiClient(AuthService);
    const response = await client.validateGuest({
      deviceId: deviceId.value ?? '',
    });
    deviceId.value = response.deviceId;
    token.value = response.token;
    saveGuestSession(response.deviceId, response.token);
  } catch {
    // Backend unreachable or auth failed — keep browsing without a guest session.
  }
}

function startGuestSessionWatch(isBackendHealthy: ReturnType<typeof useBackendHealth>['isBackendHealthy']) {
  if (sessionWatchStarted) {
    return;
  }
  sessionWatchStarted = true;

  watch(
    isBackendHealthy,
    (healthy) => {
      if (!healthy) {
        return;
      }
      if (ensureInFlight) {
        return;
      }
      ensureInFlight = ensureGuestSession(true).finally(() => {
        ensureInFlight = null;
      });
    },
    { immediate: true },
  );
}

export function useGuestAuth() {
  const { isBackendHealthy } = useBackendHealth();
  startGuestSessionWatch(isBackendHealthy);

  return {
    deviceId,
    token,
    isGuest,
  };
}
