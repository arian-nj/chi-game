package api

import (
	"context"
	"log/slog"

	"github.com/arian-nj/chigame/backend/database"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
)

type RoomMessageCommand struct {
	room    *Room
	sender  *RoomMember
	msg     *roomv1.ChatMessageRequest
	Queries *database.Queries
}

func NewRoomMessageCommand(room *Room, sender *RoomMember, msg *roomv1.ChatMessageRequest, queries *database.Queries) *RoomMessageCommand {
	return &RoomMessageCommand{
		room:    room,
		sender:  sender,
		msg:     msg,
		Queries: queries,
	}
}

func (c *RoomMessageCommand) Execute() {
	messageRow, err := c.Queries.InsertRoomMessage(context.Background(), database.InsertRoomMessageParams{
		RoomID:   c.room.ID,
		PersonID: c.sender.PersonID,
		Message:  c.msg.Text,
	})
	if err != nil {
		slog.Error("failed to insert room message", "error", err)
		return
	}

	for _, member := range c.room.Members {
		if member.PersonID == c.sender.PersonID {
			continue // Don't send to the sender
		}
		chatResponse := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Chat{
				Chat: &roomv1.ChatMessageResponse{
					PlayerId: c.sender.PersonID,
					Text:     c.msg.GetText(),
					Id:       messageRow.ID,
				},
			},
		}
		if member.Socket != nil {
			member.Socket.SendMessage(chatResponse)
		}
	}
}
