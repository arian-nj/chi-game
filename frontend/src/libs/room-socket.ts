import { getApiBaseUrl } from "./api-base-url";
import { readGuestToken } from "./guest-auth-storage";
import { create, fromBinary, toBinary } from "@bufbuild/protobuf";
import type { ChatMessageResponse } from "@/gen/room/v1/room_pb";
import { ChatMessageRequestSchema, RoomErrorType, RoomMessageSchema } from "@/gen/room/v1/room_pb";

function roomWebSocketURL(code:string): string {
  return `${getApiBaseUrl()}/room/websocket?auth_token=${readGuestToken()}&code=${code}`;
}


export function useRoomSocket(code:string) {
    const socket = new SessionSocket(roomWebSocketURL(code));
    return {
        socket,
    }
}


export class SessionSocket extends WebSocket {
    HandleChatMessage: ((chatMsg: ChatMessageResponse) => void) | null = null
  
    HandleSessionErrorMessage: ((errType: RoomErrorType) => void) | null = null
  
    // HandleChangeGametype: ((chngGameMessage: ChangeGameTypeMessage) => void) | null = null
  
    // HandleGameTimeSyncMessage: ((timeMessage: Time) => void) | null = null
  
    constructor(url: string) {
        super(url, [])
        this.binaryType = "arraybuffer"
        this.onopen = () => {
            console.log("WebSocket opened")
        }

        this.onclose = (event) => {
            console.log("WebSocket closed:", {
                code: event.code,
                reason: event.reason,
                wasClean: event.wasClean,
            });
        }

        this.onerror = (event) => {
            console.error("WebSocket error:",event);
        }

  
        this.onmessage = async (event) => {
        const bytes = new Uint8Array(event.data)
        const newSessionMessage = fromBinary(RoomMessageSchema, bytes)
        if (newSessionMessage.content.case == "chat") {
            if (this.HandleChatMessage != null) {
                this.HandleChatMessage(newSessionMessage.content.value)
            }
        } else if (newSessionMessage.content.case == "error") {
            if (this.HandleSessionErrorMessage != null) {
            this.HandleSessionErrorMessage(newSessionMessage.content.value)
            } else {
            console.error("no Session Error handler is set", newSessionMessage.content.value)
            }
        }
        //  else if (newSessionMessage.content.case == "gameType") {
        //   if (this.HandleChangeGametype != null) {
        //     this.HandleChangeGametype(newSessionMessage.content.value)
        //   } else {
        //     console.error("no Session Change Game Type handler is set", newSessionMessage.content.value)
        //   }
        // } else if (newSessionMessage.content.case == "syncTime") {
        //   if (this.HandleGameTimeSyncMessage != null) {
        //     this.HandleGameTimeSyncMessage(newSessionMessage.content.value)
        //   } else {
        //     console.error("no Session Game time sync handler is set", newSessionMessage.content.value)
        //   }
        // }
        }
  
    };
  
    SendChatReqMessage(text: string) {
      const chatReq = create(ChatMessageRequestSchema, { text });
  
      const sessionMsg = create(RoomMessageSchema, {
        content: { case: "chatReq", value: chatReq }
      });
  
      const bytes = toBinary(RoomMessageSchema, sessionMsg);
      this.send(bytes);
    }
  }