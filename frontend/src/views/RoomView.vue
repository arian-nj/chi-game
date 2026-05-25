<script setup lang="ts">
import XoOnline from '../components/game/xo/XoOnline.vue';
import Chat from '../components/chat/GameChat.vue';
import { useToast } from '../components/Toast.vue';
import { GameType, RoomErrorType } from '../gen/room/v1/room_pb';

import { GetJwtToken } from "../lib/auth";
import { GetApiUrl } from "../lib/baseURL";
import { RoomSocket } from "../lib/RoomWs";
import router from '../router/router';
import { ref, useTemplateRef, watch } from 'vue';
import Conn4Online from '../components/game/conn4/Conn4Online.vue';

const { toast } = useToast()
const isConnected = ref(false)

const activeGame = ref<null | GameType>(null)

const roomAPIUrl = GetApiUrl() + "/api/room/" + "?auth_token=" + GetJwtToken()
const roomSocket = new RoomSocket(roomAPIUrl)

const ChatRef = useTemplateRef('chat-ref')

roomSocket.onopen = () => {
  isConnected.value = true
  console.log("game room WebSocket connection established");
}

roomSocket.onerror = (event) => {
  console.error("WebSocket error:", event);
  toast.error("Socket Error")
};

function HandleError(errType: RoomErrorType) {
  if (errType == RoomErrorType.AUTH) {
    toast.error("auth error")
    router.push('/');
  } else if (errType == RoomErrorType.NOROOM) {
    toast.error("تو هیچ بازی ای نیستی")
    router.push('/');
  }
}
roomSocket.HandleRoomErrorMessage = HandleError

roomSocket.HandleChangeGametype = (changeGameMessage) => {
  const gameType = changeGameMessage.gameType
  console.log("change game type " + gameType)
  switch (gameType) {
    case GameType.XO3X3:
    case GameType.CONN4:
      activeGame.value = gameType
      break;

    default:
      return
  }
}
watch(ChatRef, () => {
  if (ChatRef.value) {
    roomSocket.HandleChatMessage = ChatRef.value.HandleIncomingChat
  }
})
</script>

<template>
  <div class="flex w-screen h-screen items-center justify-center bg-[#14bd96]">
    <div v-if="isConnected" class="w-auto h-full overflow-hidden relative flex items-center justify-center
      aspect-9/16 m-5">
      <XoOnline v-if="activeGame === GameType.XO3X3" :room-socket="roomSocket" />
      <Conn4Online v-else-if="activeGame === GameType.CONN4" :room-socket="roomSocket" />
      <h1 v-else>No Game</h1>

    </div>

    <div v-else class="text-4xl">
      Connecting...
    </div>

    <Chat :roomSocket="roomSocket" ref='chat-ref' />
  </div>
</template>
