import type { Interceptor, Transport } from '@connectrpc/connect';
import { addStaticKeyToTransport } from '@connectrpc/connect-query-core';
import { createConnectTransport } from '@connectrpc/connect-web';
import { readGuestToken } from './guest-auth-storage';
import { getApiBaseUrl } from './api-base-url';

const TRANSPORT_KEY = 'chigame-api';

const authInterceptor: Interceptor = (next) => async (request) => {
  const guestToken = readGuestToken();
  if (guestToken) {
    request.header.set('Authorization', `Bearer ${guestToken}`);
  }
  return next(request);
};

let transport: Transport | null = null;

export function getApiTransport(): Transport {
  if (!transport) {
    transport = addStaticKeyToTransport(
      createConnectTransport({
        baseUrl: getApiBaseUrl(),
        interceptors: [authInterceptor],
      }),
      TRANSPORT_KEY,
    );
  }
  return transport;
}
