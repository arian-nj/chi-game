import type { DescService } from '@bufbuild/protobuf';
import type { Interceptor, Transport } from '@connectrpc/connect';
import { createClient, type Client } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { readGuestToken } from './guest-auth-storage';
import { getApiBaseUrl } from './api-base-url';

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
    transport = createConnectTransport({
      baseUrl: getApiBaseUrl(),
      interceptors: [authInterceptor],
    });
  }
  return transport;
}

export function createApiClient<T extends DescService>(service: T): Client<T> {
  return createClient(service, getApiTransport());
}
