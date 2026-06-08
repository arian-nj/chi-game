package tictactoe

import (
	"log/slog"

	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	tictacv1 "github.com/arian-nj/chigame/backend/gen/tictac/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/arian-nj/chigame/backend/internals/room"
)

type TicState int

const (
	TicStatePlaying TicState = iota
	TicStateEnded
)

type TicEngine struct {
	board           *Board
	commander       *commander.Commander
	members         room.MapRoomMembers
	playerTurnIndex int64

	State TicState
}

func (engine *TicEngine) Init(commander *commander.Commander, members room.MapRoomMembers) {
	board, err := NewBoard(Size3)
	if err != nil {
		panic(err)
	}

	engine.board = board
	engine.members = members
	engine.commander = commander
	engine.playerTurnIndex = int64(random.GenerateRandomNumber(len(members)))
	engine.State = TicStatePlaying
}

func (engine *TicEngine) nextPlayerTurn() {
	engine.playerTurnIndex = (engine.playerTurnIndex + 1) % int64(len(engine.members))
}

func (engine *TicEngine) currentPlayerTurn() *room.RoomMember {
	return engine.members[engine.playerTurnIndex]
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
		engine.commander.PushCommand(NewPlayCommand(engine, playerId, newTicMessage.GetPlay().GetCellIndex()))
	}
}

// /
type PlayCommand struct {
	Engine    *TicEngine
	PlayerID  int64
	CellIndex int32
}

func NewPlayCommand(engine *TicEngine, playerId int64, cellIndex int32) *PlayCommand {
	return &PlayCommand{Engine: engine, PlayerID: playerId, CellIndex: cellIndex}
}

func (command *PlayCommand) Execute() {
	if command.Engine.State != TicStatePlaying {
		slog.Error("game is not in playing state", "state", command.Engine.State)
		return
	}

	currentPlayer := command.Engine.currentPlayerTurn()
	if currentPlayer.Person.ID != command.PlayerID {
		slog.Error("current player is not the player who made the move", "current_player", currentPlayer.Person.ID, "player_id", command.PlayerID)
		return
	}

	winResult, err := command.Engine.board.Play(int(command.CellIndex), Cell(currentPlayer.Person.ID))
	if err != nil {
		slog.Error("failed to play the move", "error", err)
		return
	}

	var endGameCommand *EndGameCommand = nil
	if winResult != nil {
		endGameCommand = NewEndGameCommand(command.Engine, currentPlayer, nil, tictacv1.EndReason_END_REASON_FULL)
	} else if command.Engine.board.IsFull() {
		endGameCommand = NewEndGameCommand(command.Engine, nil, nil, tictacv1.EndReason_END_REASON_TIE)
	}
	if endGameCommand != nil {
		command.Engine.State = TicStateEnded
		command.Engine.commander.InjectCommand(endGameCommand)
		return
	}

	command.Engine.nextPlayerTurn()
}

func (pCommand *PlayCommand) BroadcastNewMove() *tictacv1.PlayResponse {
	for _, member := range pCommand.Engine.members {

		if member.Person.ID == pCommand.PlayerID {
			newRoomMessage := roomv1.RoomMessage{
				Content: &roomv1.RoomMessage_Game{
					Game: &roomv1.GameMessage{
						Game: &roomv1.GameMessage_Tictactoe{
							Tictactoe: &tictacv1.TicTacToeGameMessage{
								Payload: &tictacv1.TicTacToeGameMessage_PlayResponse{
									PlayResponse: &tictacv1.PlayResponse{
										IsValid: true,
										Move: &tictacv1.Move{
											CellIndex: int32(pCommand.CellIndex),
										},
									},
								},
							},
						},
					},
				},
			}

			err := member.Socket.SendMessage(&newRoomMessage)
			if err != nil {
				slog.Error("can't send play response", "err", err)
			}
		} else {
			newRoomMessage := roomv1.RoomMessage{
				Content: &roomv1.RoomMessage_Game{
					Game: &roomv1.GameMessage{
						Game: &roomv1.GameMessage_Tictactoe{
							Tictactoe: &tictacv1.TicTacToeGameMessage{
								Payload: &tictacv1.TicTacToeGameMessage_Move{
									Move: &tictacv1.Move{
										CellIndex: int32(pCommand.CellIndex),
									},
								},
							},
						},
					},
				},
			}
			err := member.Socket.SendMessage(&newRoomMessage)
			if err != nil {
				slog.Error("can't send new move", "err", err)
			}
		}

	}
	return nil
}

// /
type StartCommand struct {
	Engine *TicEngine
}

func NewStartCommand(engine *TicEngine) *StartCommand {
	return &StartCommand{Engine: engine}
}

func (startGame *StartCommand) Execute() {
	startGame.Engine.State = TicStatePlaying
}

// /
type EndGameCommand struct {
	Reason tictacv1.EndReason
	Winner *room.RoomMember
	Loser  *room.RoomMember
	Engine *TicEngine
}

func NewEndGameCommand(engine *TicEngine, winner *room.RoomMember, loser *room.RoomMember, reason tictacv1.EndReason) *EndGameCommand {
	return &EndGameCommand{
		Reason: reason,
		Winner: winner,
		Loser:  loser,
		Engine: engine,
	}
}

func (endGame *EndGameCommand) Execute() {
}

func (endGame *EndGameCommand) BroadcastEndGame() {
	reason := tictacv1.EndReason_END_REASON_UNSPECIFIED

	newRoomMessage := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{Game: &roomv1.GameMessage_Tictactoe{Tictactoe: &tictacv1.TicTacToeGameMessage{
				Payload: &tictacv1.TicTacToeGameMessage_EndGame{
					EndGame: &tictacv1.EndGame{
						Reason: reason,
						Winner: &accountv1.Account{
							Id:       int64(endGame.Winner.Person.ID),
							Username: endGame.Winner.Person.Username,
						},
						Loser: &accountv1.Account{
							Id:       int64(endGame.Loser.Person.ID),
							Username: endGame.Loser.Person.Username,
						},
					},
				},
			}}},
		},
	}

	for _, member := range endGame.Engine.members {
		err := member.Socket.SendMessage(&newRoomMessage)
		if err != nil {
			slog.Error("error sending end game", "err", err)
		}
	}
}
