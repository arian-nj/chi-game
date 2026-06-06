import { RoomService } from '@/gen/room/v1/room_pb';
import { createApiClient } from '@/libs/api-client';

export async function createRoom(): Promise<string> {
  const client = createApiClient(RoomService);
  const response = await client.createRoom({ gameKey: '' });
  return response.code;
}

export async function joinRoomWithCode(code: string): Promise<void> {
  const client = createApiClient(RoomService);
  await client.joinRoom({
    code: code.trim().toUpperCase(),
  });
}

export async function leaveCurrentRoom(code: string): Promise<void> {
  const client = createApiClient(RoomService);
  await client.leaveRoom({
    code: code.trim().toUpperCase(),
  });
}
