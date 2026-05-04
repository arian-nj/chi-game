<script setup lang="ts">
import Conn4Board from "../../../components/game/conn4/Conn4Board.vue"
import { ref, useTemplateRef } from "vue";

import * as Conn4Buff from "../../../gen/conn4_game/v1/conn4_pb";
import { create, toBinary } from "@bufbuild/protobuf";
import { useToast } from "../../../components/Toast.vue";
import EndGame from "../EndGame.vue";
import PlayersBoard from "../PlayersBoard.vue";
import { RoomMessageSchema } from "../../../gen/room/v1/room_pb";
import type { RoomSocket } from "../../../lib/RoomWs";


const props = defineProps({
  roomSocket: {
    type: Object as () => RoomSocket,
    required: true
  }
})

const Conn4BoardRef = useTemplateRef('conn4-board-ref')
const EndGameData = ref<Conn4Buff.EndGame>()

const { toast } = useToast()

props.roomSocket.HandleGameMessage = (gameMessage) => {
  console.log(gameMessage)
  if (gameMessage.game.case != "conn4") {
    throw Error("non xo message ended up in xo")
  }
  const payload = gameMessage.game.value.payload
  switch (payload.case) {
    case "gameState":
      Conn4BoardRef.value?.SetCells(payload.value.cells)
      break

    case "move":
      handleMoveAction(payload.value)
      break
    case "playResponse":
      handlePlayResponse(payload.value)
      break
    // case "syncTime":
    //   PlayersBoardRef.value?.handleTimeSync(payload.value)
    //   break
    case "endGame":
      EndGameData.value = payload.value
  }
}
const handleMoveAction = (moveData: Conn4Buff.Move) => {
  Conn4BoardRef.value?.DoMove(moveData.cellIndex, moveData.cellValue)
}

function sendClick(i: number) {
  console.log(i, " clicked")
  const newRoomMsg = create(RoomMessageSchema, {
    content: {
      case: "game",
      value: {
        game: {
          case: "conn4",
          value: { payload: { case: "play", value: { colIndex: i } } }
        }
      }
    }
  })
  const bytes = toBinary(RoomMessageSchema, newRoomMsg)
  props.roomSocket.send(bytes)
}

const handlePlayResponse = (playResponse: Conn4Buff.PlayResponse) => {
  if (playResponse.isValid) {
    if (playResponse.move) {
      Conn4BoardRef.value?.DoMove(playResponse.move.cellIndex, playResponse.move.cellValue)
    }
  } else {
    toast.error(playResponse.reason)
  }
}


</script>

<template>
  <PlayersBoard :roomSocket="roomSocket" />
  <Conn4Board ref="conn4-board-ref" @col-selected="sendClick" />
  <EndGame v-if="EndGameData" :loser="EndGameData.loser" :winner="EndGameData.winner" />
</template>
