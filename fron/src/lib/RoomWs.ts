import { create, fromBinary, toBinary } from "@bufbuild/protobuf";
import { ChatMessageRequestSchema, RoomMessageSchema, type ChangeGameTypeMessage, type ChatMessage, type GameMessage, type RoomErrorType, type Time } from "../gen/room/v1/room_pb";

export class RoomSocket extends WebSocket {
  HandleChatMessage: ((chatMsg: ChatMessage) => void) | null = null

  HandleGameMessage: ((msg: GameMessage) => void) | null = null

  HandleRoomErrorMessage: ((errType: RoomErrorType) => void) | null = null

  HandleChangeGametype: ((chngGameMessage: ChangeGameTypeMessage) => void) | null = null

  HandleGameTimeSyncMessage: ((timeMessage: Time) => void) | null = null

  constructor(url: string) {
    super(url, [])
    this.binaryType = "arraybuffer"
    this.onclose = (event) => {
      console.warn("WebSocket closed:", {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean,
      });
    };

    this.onmessage = async (event) => {
      const bytes = new Uint8Array(event.data)
      const newRoomMessage = fromBinary(RoomMessageSchema, bytes)
      if (newRoomMessage.content.case == "chat") {
        if (this.HandleChatMessage != null) {
          this.HandleChatMessage?.(newRoomMessage.content.value)
        } else {
          throw new Error("no chat handler is set")
        }
      } else if (newRoomMessage.content.case == "game") {
        if (this.HandleGameMessage != null) {
          this.HandleGameMessage(newRoomMessage.content.value)
        } else {
          console.error("no game message handler is set", newRoomMessage.content.value)
        }
      } else if (newRoomMessage.content.case == "error") {
        if (this.HandleRoomErrorMessage != null) {
          this.HandleRoomErrorMessage(newRoomMessage.content.value)
        } else {
          console.error("no Room Error handler is set", newRoomMessage.content.value)
        }
      } else if (newRoomMessage.content.case == "gameType") {
        if (this.HandleChangeGametype != null) {
          this.HandleChangeGametype(newRoomMessage.content.value)
        } else {
          console.error("no Room Change Game Type handler is set", newRoomMessage.content.value)
        }
      } else if (newRoomMessage.content.case == "syncTime") {
        if (this.HandleGameTimeSyncMessage != null) {
          this.HandleGameTimeSyncMessage(newRoomMessage.content.value)
        } else {
          console.error("no Room Game time sync handler is set", newRoomMessage.content.value)
        }
      }


    }

  };

  SendChatReqMessage(text: string) {
    const chatReq = create(ChatMessageRequestSchema, { text });

    const RoomMsg = create(RoomMessageSchema, {
      content: { case: "chatReq", value: chatReq }
    });

    const bytes = toBinary(RoomMessageSchema, RoomMsg);
    this.send(bytes);
  }
}


