package room

import (
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

// Game Engine use room member
type GameEngine interface {
	Init(commaner *commander.Commander, roomMembers MapRoomMembers)

	Key() string
	MinPlayers() int
	MaxPlayers() int

	SocketRouter(gameMessage *roomv1.GameMessage, playerId int64)
}
