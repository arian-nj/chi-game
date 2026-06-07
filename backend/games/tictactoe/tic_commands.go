package tictactoe

import (
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	tictacv1 "github.com/arian-nj/chigame/backend/gen/tictac/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

type TicEngine struct {
	board     *Board
	commander *commander.Commander
}

func (engine *TicEngine) Init(commander *commander.Commander) {
	engine.commander = commander
	board, err := NewBoard(Size3)
	if err != nil {
		panic(err)
	}
	engine.board = board
}

func (engine *TicEngine) Key() string {
	return "tic-tac-toe"
}

func (engine *TicEngine) MinPlayers() int {
	return 2
}

func (engine *TicEngine) MaxPlayers() int {
	return 2
}

func (engine *TicEngine) SocketRouter(gameMessage *roomv1.GameMessage, playerId int64) {
	newTicMessage := gameMessage.GetTictactoe()
	switch newTicMessage.Payload.(type) {
	case *tictacv1.TicTacToeGameMessage_Play:
		engine.commander.PushCommand(NewPlayCommand(engine.board, playerId, newTicMessage.GetPlay().GetCellIndex()))
	}
}

type PlayCommand struct {
	board     *Board
	playerId  int64
	cellIndex int32
}

func NewPlayCommand(board *Board, playerId int64, cellIndex int32) *PlayCommand {
	return &PlayCommand{board: board, playerId: playerId, cellIndex: cellIndex}
}

func (command *PlayCommand) Execute() {
}
