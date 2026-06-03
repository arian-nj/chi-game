package conn4

import (
	"log/slog"

	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	conn4_gamev1 "github.com/arian-nj/chigame/backend/gen/conn4_game/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

func (game *Conn4State) SocketRouter(newGameMsg *roomv1.GameMessage, playerId int) {

	newXoMessage := newGameMsg.GetConn4()
	switch newXoMessage.Payload.(type) {
	case *conn4_gamev1.Conn4GameMessage_Play:
		playData := newXoMessage.GetPlay()
		PlayHandlerSocket(game, playData, playerId)
	}
}

type SocketListener struct {
	Player *Conn4Player
}

func (sl *SocketListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MoveCommand:
		sl.SocketBrodcastNewMove(c)
	// case *StartCommand:
	case *EndGameCommand:
		sl.SendEndGameSocket(c)
	case *SyncTimeCommand:
		sl.SocketBrodcastSyncTime(c.Game)
	}
}

func sendInvalidResponse(player *Conn4Player, errMsg string) error {
	newRoomMsg := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{
				Game: &roomv1.GameMessage_Conn4{
					Conn4: &conn4_gamev1.Conn4GameMessage{
						Payload: &conn4_gamev1.Conn4GameMessage_PlayResponse{
							PlayResponse: &conn4_gamev1.PlayResponse{
								IsValid: false,
								Reason:  errMsg,
							},
						},
					},
				},
			},
		},
	}

	return player.Socket.SendMessage(&newRoomMsg)
}

func (sl *SocketListener) SocketBrodcastSyncTime(gameState *Conn4State) {
	allTimeMessages := []*roomv1.RoomMessage{}
	for _, player := range gameState.Players {
		newRoomMessage := roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_SyncTime{
				SyncTime: &roomv1.Time{
					PlayerId:  int64(player.ID),
					SpentTime: int32(player.Timer.SpentInt()),
					TotalTime: int32(MAX_ALLOWED_TIME_INT),
				},
			},
		}
		allTimeMessages = append(allTimeMessages, &newRoomMessage)
	}

	player := gameState.findByID(sl.Player.ID)
	for _, timeEvent := range allTimeMessages {
		err := player.Socket.SendMessage(timeEvent)
		if err != nil {
			slog.Error("can't send new move", "err", err)
		}
	}
}

func SocketSendGameState(gameState *Conn4State, player *Conn4Player) {
	cells := []int32{}
	for _, cell := range gameState.Board.Board {
		cells = append(cells, int32(cell))
	}
	newRoomMessage := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{
				Game: &roomv1.GameMessage_Conn4{
					Conn4: &conn4_gamev1.Conn4GameMessage{
						Payload: &conn4_gamev1.Conn4GameMessage_GameState{
							GameState: &conn4_gamev1.GameState{
								Cells: cells,
							},
						},
					},
				},
			},
		},
	}

	err := player.Socket.SendMessage(&newRoomMessage)
	if err != nil {
		slog.Error("error sending game state")
	}
}

func (sl *SocketListener) SocketBrodcastNewMove(moveCommand *MoveCommand) {

	if sl.Player.ID == moveCommand.PlayerID {
		newRoomMessage := roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Game{
				Game: &roomv1.GameMessage{
					Game: &roomv1.GameMessage_Conn4{
						Conn4: &conn4_gamev1.Conn4GameMessage{
							Payload: &conn4_gamev1.Conn4GameMessage_PlayResponse{
								PlayResponse: &conn4_gamev1.PlayResponse{
									IsValid: true,
									Move: &conn4_gamev1.Move{
										CellIndex: int32(moveCommand.CellIndex),
										CellValue: int32(moveCommand.MoveType),
									},
								},
							},
						},
					},
				},
			},
		}

		err := sl.Player.Socket.SendMessage(&newRoomMessage)
		if err != nil {
			slog.Error("can't send play response", "err", err)
		}
	} else {
		newRoomMessage := roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Game{
				Game: &roomv1.GameMessage{
					Game: &roomv1.GameMessage_Conn4{
						Conn4: &conn4_gamev1.Conn4GameMessage{
							Payload: &conn4_gamev1.Conn4GameMessage_Move{
								Move: &conn4_gamev1.Move{
									CellValue: int32(moveCommand.MoveType),
									CellIndex: int32(moveCommand.CellIndex),
								},
							},
						},
					},
				},
			},
		}
		err := sl.Player.Socket.SendMessage(&newRoomMessage)
		if err != nil {
			slog.Error("can't send new move", "err", err)
		}
	}
}

func (sl *SocketListener) SendEndGameSocket(endGameCommand *EndGameCommand) {
	reason := conn4_gamev1.EndReason_END_REASON_UNSPECIFIED
	switch endGameCommand.reason {
	case END_GAME_TIE:
		reason = conn4_gamev1.EndReason_END_REASON_TIE
	case END_GAME_FULL:
		reason = conn4_gamev1.EndReason_END_REASON_FULL
	case END_GAME_TIMEOUT:
		reason = conn4_gamev1.EndReason_END_REASON_TIMOUT

	}

	newRoomMessage := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{
				Game: &roomv1.GameMessage_Conn4{
					Conn4: &conn4_gamev1.Conn4GameMessage{
						Payload: &conn4_gamev1.Conn4GameMessage_EndGame{
							EndGame: &conn4_gamev1.EndGame{
								Reason: reason,
								Winner: &accountv1.Account{
									Id:       int64(endGameCommand.Winner.ID),
									Username: endGameCommand.Winner.Name,
								},
								Loser: &accountv1.Account{
									Id:       int64(endGameCommand.Loser.ID),
									Username: endGameCommand.Loser.Name,
								},
							},
						},
					},
				},
			},
		},
	}
	err := sl.Player.Socket.SendMessage(&newRoomMessage)
	if err != nil {
		slog.Error("error sending game state")
	}

}

// routed

func PlayHandlerSocket(game *Conn4State, playInput *conn4_gamev1.Play, playerID int) {
	player := game.findByID(playerID)

	if player == nil {
		slog.Error("can't find player in socjet play handler")
		return
	}

	if game.CurrentPlayer().ID != playerID {
		sendInvalidResponse(player,
			"نوبت تو نیست")
		return
	}

	moveType := game.CurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(int(playInput.ColIndex))

	if !isValid {
		err := sendInvalidResponse(player, errMsg)
		if err != nil {
			slog.Error("can't send invalid response")
		}
		return
	}

	playCommand := NewMoveCommand(game, int(playInput.ColIndex), moveType, player.ID)
	game.PushCommand(playCommand)
}
