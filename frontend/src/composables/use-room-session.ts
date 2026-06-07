import { useToast } from '@/components/Toast.vue';
import { useGuestAuth } from '@/composables/use-guest-auth';
import { useGuestProfile } from '@/composables/use-guest-profile';
import { useTextDirection } from '@/composables/use-text-direction';
import type { Account } from '@/gen/account/v1/account_pb';
import { RoomErrorType, RoomService } from '@/gen/room/v1/room_pb';
import { createApiClient } from '@/libs/api-client';
import { joinRoomWithCode, leaveCurrentRoom } from '@/libs/room-api';
import type { SessionSocket } from '@/libs/room-socket';
import { useRoomSocket } from '@/libs/room-socket';
import { roomInviteUrl } from '@/libs/room-url';
import { ConnectError } from '@connectrpc/connect';
import { useQuery } from '@tanstack/vue-query';
import type { ComputedRef, InjectionKey, Ref } from 'vue';
import { computed, inject, onBeforeUnmount, ref, watch } from 'vue';
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useRoomChat, type RoomChatMessage } from '@/composables/use-room-chat';

export type RoomSessionContext = {
  locale: ComputedRef<string>;
  roomCode: ComputedRef<string>;
  isBusy: Ref<boolean>;
  players: ComputedRef<Account[]>;
  isError: Ref<boolean>;
  error: Ref<Error | null>;
  isReadyToPlay: ComputedRef<boolean>;
  isHost: ComputedRef<boolean>;
  roomLink: ComputedRef<string>;
  textDir: ComputedRef<'ltr' | 'rtl'>;
  socket: SessionSocket;
  chatMessages: Ref<RoomChatMessage[]>;
  chatConnected: Ref<boolean>;
  chatHasMessages: ComputedRef<boolean>;
  sendChatMessage: (text: string) => boolean;
  chatPlayerName: (playerId: bigint) => string;
  isOwnChatMessage: (playerId: bigint) => boolean;
  leave: () => Promise<void>;
  playerLabel: (displayName: string, username: string) => string;
  copyText: (text: string) => Promise<void>;
};

export const roomSessionKey: InjectionKey<RoomSessionContext> = Symbol('roomSession');

export function useRoomSession(): RoomSessionContext {
  const session = inject(roomSessionKey);
  if (!session) {
    throw new Error('useRoomSession must be used within RoomShell');
  }
  return session;
}

function inviteErrorMessage(err: unknown, t: (key: string) => string): string {
  if (err instanceof ConnectError) {
    const msg = err.message.toLowerCase();
    if (msg.includes('full')) {
      return t('invite.roomFull');
    }
    if (msg.includes('invalid') || msg.includes('expired') || msg.includes('not found')) {
      return t('invite.invalidCode');
    }
    if (err.code === 16) {
      return t('invite.needAuth');
    }
  }
  return '';
}

export function provideRoomSession(): RoomSessionContext {
  const route = useRoute();
  const router = useRouter();
  const { t } = useI18n();
  const { textDir } = useTextDirection();
  const { isGuest } = useGuestAuth();
  const { toast } = useToast();
  const { data: meData } = useGuestProfile();

  const locale = computed(() => route.params.locale as string);
  const roomCode = computed(() => {
    const raw = route.params.code;
    if (typeof raw !== 'string') {
      return '';
    }
    return raw.trim().toUpperCase();
  });

  const joinAttempted = ref(false);
  const joinReady = ref(false);
  const isBusy = ref(false);

  const client = createApiClient(RoomService);

  const { data, isError, error } = useQuery({
    queryKey: computed(() => ['room', roomCode.value]),
    queryFn: ({ signal }) => client.getRoom({ code: roomCode.value }, { signal }),
    enabled: computed(() => Boolean(roomCode.value) && joinReady.value),
  });

  const playersList = ref<Account[]>([]);
  const players = computed(() => playersList.value);

  function addPlayer(player: Account) {
    if (playersList.value.some((entry) => entry.id === player.id)) {
      return;
    }
    playersList.value = [...playersList.value, player];
  }

  function removePlayer(playerId: bigint) {
    playersList.value = playersList.value.filter((entry) => entry.id !== playerId);
  }

  watch(
    () => data.value?.hostPlayer,
    (hostPlayer) => {
      if (hostPlayer) {
        addPlayer(hostPlayer);
      }
    },
    { immediate: true },
  );
  const isReadyToPlay = computed(() => players.value.length >= 2);
  const isHost = computed(() => data.value?.hostPlayer?.id === meData.value?.account?.id);
  const roomLink = computed(() => roomInviteUrl(locale.value, roomCode.value));

  function playerLabel(displayName: string, username: string) {
    const name = displayName?.trim();
    if (name) {
      return name;
    }
    return username ? `@${username}` : t('invite.unknownPlayer');
  }

  async function copyText(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      toast.success(t('invite.copied'));
    } catch {
      // ignore
    }
  }

  async function joinRoom(code: string) {
    if (!isGuest.value) {
      toast.info(t('invite.needAuth'));
      await router.replace({ name: 'room', params: { locale: locale.value } });
      return;
    }

    const normalized = code.trim().toUpperCase();
    if (!normalized) {
      return;
    }

    isBusy.value = true;
    try {
      await joinRoomWithCode(normalized);
      joinReady.value = true;
      if (route.params.code !== normalized) {
        await router.replace({
          name: 'room-code',
          params: { locale: locale.value, code: normalized },
        });
      }
    } catch (err) {
      const specific = inviteErrorMessage(err, t);
      toast.error(specific || t('invite.joinFailed'));
      await router.push({ name: 'room', params: { locale: locale.value } });
    } finally {
      isBusy.value = false;
    }
  }

  watch(
    roomCode,
    async (code) => {
      if (!code || joinAttempted.value) {
        return;
      }
      joinAttempted.value = true;
      await joinRoom(code);
    },
    { immediate: true },
  );

  watch(
    players,
    (list) => {
      if (!roomCode.value) {
        return;
      }
      if (list.length < 2 && route.name === 'room-play') {
        void router.replace({
          name: 'room-code',
          params: { locale: locale.value, code: roomCode.value },
        });
      }
    },
    { flush: 'post' },
  );

  async function leave() {
    try {
      await leaveCurrentRoom(roomCode.value);
    } catch {
      // still leave UI
    }
    await router.push({ name: 'room', params: { locale: locale.value } });
  }

  onBeforeRouteLeave((to) => {
    if (to.name === 'room' || to.name === 'room-code' || to.name === 'room-play') {
      return;
    }
    void leaveCurrentRoom(roomCode.value);
  });

  const { socket } = useRoomSocket(roomCode.value);
  socket.HandleSessionErrorMessage = (errType: RoomErrorType) => {
    toast.error(`${errType}`);
  };
  socket.HandleMemberJoined = (memberJoined) => {
    const player = memberJoined.player;
    if (player) {
      addPlayer(player);
    }
  };
  socket.HandleMemberLeft = (memberLeft) => {
    const player = memberLeft.player;
    if (player?.id) {
      removePlayer(player.id);
    }
  };

  const myId = computed(() => meData.value?.account?.id);
  const {
    messages: chatMessages,
    isConnected: chatConnected,
    hasMessages: chatHasMessages,
    sendMessage: sendChatMessage,
    playerNameFor: chatPlayerName,
    isOwnMessage: isOwnChatMessage,
  } = useRoomChat(socket, players, myId, playerLabel);

  onBeforeUnmount(() => {
    socket.HandleMemberJoined = null;
    socket.HandleMemberLeft = null;
    socket.close();
  });

  return {
    locale,
    roomCode,
    isBusy,
    players,
    isError,
    error,
    isReadyToPlay,
    isHost,
    roomLink,
    textDir,
    socket,
    chatMessages,
    chatConnected,
    chatHasMessages,
    sendChatMessage,
    chatPlayerName,
    isOwnChatMessage,
    leave,
    playerLabel,
    copyText,
  };
}
