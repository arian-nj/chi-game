import type { Account } from '@/gen/account/v1/account_pb';
import type { ChatMessageResponse } from '@/gen/room/v1/room_pb';
import type { SessionSocket } from '@/libs/room-socket';
import type { ComputedRef, Ref } from 'vue';
import { computed, onScopeDispose, ref } from 'vue';

export type RoomChatMessage = {
  id: string;
  playerId: bigint;
  text: string;
};

let serverFallbackId = 0;

function chatMessageId(chatMsg: ChatMessageResponse): string {
  if (chatMsg.id !== 0n) {
    return `srv-${chatMsg.id}`;
  }
  // Server should assign ids; this fallback must not key on text (duplicate messages are valid).
  return `srv-fallback-${serverFallbackId++}`;
}

export function useRoomChat(
  socket: SessionSocket,
  players: ComputedRef<Account[]>,
  myId: ComputedRef<bigint | undefined>,
  playerLabel: (displayName: string, username: string) => string,
) {
  const messages = ref<RoomChatMessage[]>([]);
  const isConnected = ref(socket.readyState === WebSocket.OPEN);
  let localId = 0;

  function onOpen() {
    isConnected.value = true;
  }

  function onClose() {
    isConnected.value = false;
  }

  function onChat(chatMsg: ChatMessageResponse) {
    const id = chatMessageId(chatMsg);
    if (messages.value.some((message) => message.id === id)) {
      return;
    }
    messages.value.push({
      id,
      playerId: chatMsg.playerId,
      text: chatMsg.text,
    });
  }

  socket.HandleChatMessage = onChat;
  socket.addEventListener('open', onOpen);
  socket.addEventListener('close', onClose);

  onScopeDispose(() => {
    if (socket.HandleChatMessage === onChat) {
      socket.HandleChatMessage = null;
    }
    socket.removeEventListener('open', onOpen);
    socket.removeEventListener('close', onClose);
  });

  function playerNameFor(playerId: bigint): string {
    const player = players.value.find((entry) => entry.id === playerId);
    if (!player) {
      return playerLabel('', '');
    }
    return playerLabel(player.displayName, player.username);
  }

  function isOwnMessage(playerId: bigint): boolean {
    const id = myId.value;
    return id !== undefined && playerId === id;
  }

  function sendMessage(text: string): boolean {
    const trimmed = text.trim();
    if (!trimmed || socket.readyState !== WebSocket.OPEN) {
      return false;
    }

    socket.SendChatReqMessage(trimmed);

    const id = myId.value;
    if (id !== undefined) {
      messages.value.push({
        id: `local-${localId++}`,
        playerId: id,
        text: trimmed,
      });
    }

    return true;
  }

  const hasMessages = computed(() => messages.value.length > 0);

  return {
    messages: messages as Ref<RoomChatMessage[]>,
    isConnected,
    hasMessages,
    sendMessage,
    playerNameFor,
    isOwnMessage,
  };
}
