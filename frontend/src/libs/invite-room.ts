import { InviteService } from '@/gen/invite/v1/invite_pb';
import { createApiClient } from '@/libs/api-client';

export async function createInviteRoom(): Promise<string> {
  const client = createApiClient(InviteService);
  const response = await client.createInviteRoom({ gameKey: '' });
  return response.inviteCode;
}

export async function joinRoomWithCode(inviteCode: string): Promise<void> {
  const client = createApiClient(InviteService);
  await client.joinInviteRoom({
    inviteCode: inviteCode.trim().toUpperCase(),
  });
}

export async function leaveCurrentRoom(): Promise<void> {
  const client = createApiClient(InviteService);
  await client.leaveInviteRoom({});
}
