package xo

import (
	"log/slog"

	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	xo_gamev1 "github.com/arian-nj/chigame/backend/gen/xo_game/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

// FIXME: make it map handler with auto Converting inputs
func (game *XOState) SocketRouter(newGameMsg *roomv1.GameMessage, playerId int) {

	newXoMessage := newGameMsg.GetXo()
	switch newXoMessage.Payload.(type) {
	case *xo_gamev1.XoGameMessage_Play:
		playData := newXoMessage.GetPlay()
		PlayHandlerSocket(game, playData, playerId)
	}
}

type SocketListener struct {
	Player *XoPlayer
}

func (sl *SocketListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MoveCommand:
		sl.SocketBrodcastNewMove(c.Game, c.Pos, c.MoveType, c.PlayerID)
	case *StartCommand:
	case *EndGameCommand:
		sl.SendEndGameSocket(c)
	case *SyncTimeCommand:
		sl.SocketBrodcastSyncTime(c.Game)
	}
}

func sendInvalidResponse(player *XoPlayer, errMsg string, cellIndex int32, cellType Cell) error {
	newRoomMsg := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{
				Game: &roomv1.GameMessage_Xo{
					Xo: &xo_gamev1.XoGameMessage{
						Payload: &xo_gamev1.XoGameMessage_PlayResponse{
							PlayResponse: &xo_gamev1.PlayResponse{
								IsValid: false,
								Reason:  errMsg,
								Move: &xo_gamev1.Move{
									CellIndex: cellIndex,
									CellValue: int32(cellType),
								},
							},
						},
					},
				},
			},
		},
	}

	return player.Socket.SendMessage(&newRoomMsg)
}

func PlayHandlerSocket(game *XOState, playInput *xo_gamev1.Play, playerID int) {
	player := game.findByID(playerID)

	if player == nil {
		slog.Error("can't find player in socjet play handler")
		return
	}

	if game.CurrentPlayer().ID != playerID {
		sendInvalidResponse(player, "نوبت تو نیست", playInput.CellIndex, 0)
		return
	}

	moveType := game.CurrentPlayer().Move

	isValid, errMsg := game.Board.IsMoveValid(int(playInput.CellIndex), moveType)

	if !isValid {
		err := sendInvalidResponse(player, errMsg, playInput.CellIndex, 0)
		if err != nil {
			slog.Error("can't send invalid response")
		}
		return
	}

	playCommand := NewMoveCommand(game, int(playInput.CellIndex), moveType, player.ID)
	game.PushCommand(playCommand)
}

func (sl *SocketListener) SocketBrodcastNewMove(game *XOState, moveIndex int, cellType Cell, playerID int) {

	if sl.Player.ID == playerID {
		newRoomMessage := roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Game{
				Game: &roomv1.GameMessage{
					Game: &roomv1.GameMessage_Xo{
						Xo: &xo_gamev1.XoGameMessage{
							Payload: &xo_gamev1.XoGameMessage_PlayResponse{
								PlayResponse: &xo_gamev1.PlayResponse{
									IsValid: true,
									Move: &xo_gamev1.Move{
										CellIndex: int32(moveIndex),
										CellValue: int32(cellType),
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
					Game: &roomv1.GameMessage_Xo{
						Xo: &xo_gamev1.XoGameMessage{
							Payload: &xo_gamev1.XoGameMessage_Move{
								Move: &xo_gamev1.Move{
									CellIndex: int32(moveIndex),
									CellValue: int32(cellType),
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

func (sl *SocketListener) SocketBrodcastSyncTime(gameState *XOState) {
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

func (sl *SocketListener) SendEndGameSocket(endGameCommand *EndGameCommand) {
	reason := xo_gamev1.EndReason_END_REASON_UNSPECIFIED
	switch endGameCommand.reason {
	case END_GAME_TIE:
		reason = xo_gamev1.EndReason_END_REASON_TIE
	case END_GAME_FULL:
		reason = xo_gamev1.EndReason_END_REASON_FULL
	case END_GAME_TIMEOUT:
		reason = xo_gamev1.EndReason_END_REASON_TIMOUT

	}

	newRoomMessage := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{Game: &roomv1.GameMessage_Xo{Xo: &xo_gamev1.XoGameMessage{
				Payload: &xo_gamev1.XoGameMessage_EndGame{
					EndGame: &xo_gamev1.EndGame{
						Reason: reason,
						Winner: &accountv1.Account{
							Id:   int64(endGameCommand.Winner.ID),
							Name: endGameCommand.Winner.Name,
						},
						Loser: &accountv1.Account{
							Id:   int64(endGameCommand.Loser.ID),
							Name: endGameCommand.Loser.Name,
						},
					},
				},
			}}},
		},
	}
	err := sl.Player.Socket.SendMessage(&newRoomMessage)
	if err != nil {
		slog.Error("error sending game state")
	}

}

func SocketSendGameState(gameState *XOState, player *XoPlayer) {
	cells := []int32{}
	for _, cell := range gameState.Board.Board {
		cells = append(cells, int32(cell))
	}
	newRoomMessage := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Game{
			Game: &roomv1.GameMessage{
				Game: &roomv1.GameMessage_Xo{
					Xo: &xo_gamev1.XoGameMessage{
						Payload: &xo_gamev1.XoGameMessage_GameState{
							GameState: &xo_gamev1.GameState{
								Cells:        cells,
								TurnPlayerId: int64(gameState.CurrentPlayer().ID),
								BoardSize:    int32(gameState.Board.MaxCellSize),
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
