import { RoomService } from '@/gen/room/v1/room_pb';
import { createApiClient } from '@/libs/api-client';

export async function createRoom(): Promise<string> {
  const client = createApiClient(RoomService);
  const response = await client.createRoom({ gameKey: '' });
  return response.code;
}
