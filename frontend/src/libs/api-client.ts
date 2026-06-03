import type { Transport } from '@connectrpc/connect';
import { addStaticKeyToTransport } from '@connectrpc/connect-query-core';
import { createConnectTransport } from '@connectrpc/connect-web';
import { getApiBaseUrl } from './api-base-url';

const TRANSPORT_KEY = 'chigame-api';

let transport: Transport | null = null;

export function getApiTransport(): Transport {
  if (!transport) {
    transport = addStaticKeyToTransport(
      createConnectTransport({
        baseUrl: getApiBaseUrl(),
      }),
      TRANSPORT_KEY,
    );
  }
  return transport;
}
