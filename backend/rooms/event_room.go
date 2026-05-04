package rooms

import roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"

type RoomEvent struct {
	Player *RoomPlayer
	Event  *roomv1.RoomMessage
}

func NewRoomEvent(player *RoomPlayer, event *roomv1.RoomMessage) *RoomEvent {
	return &RoomEvent{
		Player: player,
		Event:  event,
	}
}
