package games

import (
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

type GameEngine interface {
	Init(commaner *commander.Commander)

	Key() string
	MinPlayers() int
	MaxPlayers() int

	SocketRouter(gameMessage *roomv1.GameMessage, playerId int64)
}

// type GameEngineFactory func(room *api.Room) GameEngine

// var AllGameEngines = map[string]GameEngineFactory{
// 	"tic-tac-toe": NewTicTacToeGameEngine,
// 	"connect-4":   NewConnect4GameEngine,
// }

// func GetGameEngine(gameKey string) GameEngine {
// 	factory, ok := AllGameEngines[gameKey]
// 	if !ok {
// 		return nil
// 	}
// 	return factory(room)
// }
