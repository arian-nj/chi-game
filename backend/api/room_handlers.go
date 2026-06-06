package api

import roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"

func (r *Room) SocketRequestSendMsg(sender *RoomMember, msg *roomv1.ChatMessageRequest) {
	// Deliver the chat message to all other players in the room except the sender
	for id, member := range r.Members {
		if id == sender.PersonID {
			continue // Don't send to the sender
		}
		chatResponse := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Chat{
				Chat: &roomv1.ChatMessageResponse{
					PlayerId: sender.PersonID,
					Text:     msg.GetText(),
					Id:       0, // Optionally generate a message ID here if needed
				},
			},
		}
		if member.Socket != nil {
			member.Socket.SendMessage(chatResponse)
		}
	}
}
