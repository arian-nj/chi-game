package games

type GameEngine interface {
	Key() string
	MinPlayers() int
	MaxPlayers() int
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
