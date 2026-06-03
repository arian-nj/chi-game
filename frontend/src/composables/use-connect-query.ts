import type { DescMessage, DescMethodUnary, MessageInitShape, MessageShape } from '@bufbuild/protobuf';
import {
  callUnaryMethod,
  createConnectQueryKey,
} from '@connectrpc/connect-query-core';
import { getApiTransport } from '@/libs/api-client';
import { useQuery, type QueryKey } from '@tanstack/vue-query';
import type { MaybeRefOrGetter } from 'vue';

type ConnectQueryExtraOptions = {
  enabled?: MaybeRefOrGetter<boolean>;
  refetchInterval?: number | false;
  staleTime?: number;
  retry?: number | false;
};

export function useConnectQuery<I extends DescMessage, O extends DescMessage>(
  schema: DescMethodUnary<I, O>,
  input?: MessageInitShape<I>,
  options?: ConnectQueryExtraOptions,
) {
  const transport = getApiTransport();

  return useQuery<MessageShape<O>, Error>({
    queryKey: createConnectQueryKey({
      schema,
      input,
      transport,
      cardinality: 'finite',
    }) as QueryKey,
    queryFn: ({ signal }) => callUnaryMethod(transport, schema, input, { signal }),
    ...options,
  });
}
