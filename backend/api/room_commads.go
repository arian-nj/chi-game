package api

import (
	"context"
	"log/slog"

	"github.com/arian-nj/chigame/backend/database"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
)

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
	messageRow, err := c.room.Queries.InsertRoomMessage(context.Background(), database.InsertRoomMessageParams{
		RoomID:   c.room.ID,
		PersonID: c.sender.Person.ID,
		Message:  c.msg.Text,
	})
	if err != nil {
		slog.Error("failed to insert room message", "error", err)
		return
	}

	for _, member := range c.room.Members {
		if member.Person.ID == c.sender.Person.ID {
			continue // Don't send to the sender
		}
		chatResponse := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Chat{
				Chat: &roomv1.ChatMessageResponse{
					PlayerId: c.sender.Person.ID,
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

type RoomMemberJoinedCommand struct {
	room      *Room
	newMember *RoomMember
	Queries   *database.Queries
}

func NewRoomMemberJoinedCommand(room *Room, member *RoomMember) *RoomMemberJoinedCommand {
	return &RoomMemberJoinedCommand{
		room:      room,
		newMember: member,
	}
}

func (c *RoomMemberJoinedCommand) Execute() {
	for _, member := range c.room.Members {
		person, err := c.room.Queries.GetPersonByID(context.Background(), c.newMember.Person.ID)
		if err != nil {
			slog.Error("failed to get person by id", "error", err)
			continue
		}
		account := personToAccount(&person)
		chatResponse := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_MemberJoined{
				MemberJoined: &roomv1.RoomMemberJoined{
					Player: account,
				},
			},
		}
		if member.Socket != nil {
			member.Socket.SendMessage(chatResponse)
		}
	}
}

type RoomMemberLeftCommand struct {
	room      *Room
	oldMember *RoomMember
}

func NewRoomMemberLeftCommand(room *Room, member *RoomMember) *RoomMemberLeftCommand {
	return &RoomMemberLeftCommand{
		room:      room,
		oldMember: member,
	}
}
func (c *RoomMemberLeftCommand) Execute() {
	account := personToAccount(c.oldMember.Person)
	for _, member := range c.room.Members {
		if member.Person.ID == c.oldMember.Person.ID {
			continue
		}
		msg := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_MemberLeft{
				MemberLeft: &roomv1.RoomMemberLeft{
					Player: account,
				},
			},
		}
		if member.Socket != nil {
			member.Socket.SendMessage(msg)
		}
	}
}
