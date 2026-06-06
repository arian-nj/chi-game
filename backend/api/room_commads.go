package api

import roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"

type RoomMessageCommand struct {
	room   *Room
	sender *RoomMember
	msg    *roomv1.ChatMessageRequest
}

func NewRoomMessageCommand(room *Room, sender *RoomMember, msg *roomv1.ChatMessageRequest) *RoomMessageCommand {
	return &RoomMessageCommand{
		room:   room,
		sender: sender,
		msg:    msg,
	}
}

func (c *RoomMessageCommand) Execute() {
	for _, member := range c.room.Members {
		if member.PersonID == c.sender.PersonID {
			continue // Don't send to the sender
		}
		chatResponse := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Chat{
				Chat: &roomv1.ChatMessageResponse{
					PlayerId: c.sender.PersonID,
					Text:     c.msg.GetText(),
					Id:       0, // Optionally generate a message ID here if needed
				},
			},
		}
		if member.Socket != nil {
			member.Socket.SendMessage(chatResponse)
		}
	}
}
