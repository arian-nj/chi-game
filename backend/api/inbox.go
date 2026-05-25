package api

// type Inbox struct {
// 	SocketClient *socket.Socket
// }

// func NewInbox(socketClient *socket.Socket) *Inbox {
// 	return &Inbox{
// 		SocketClient: socketClient,
// 	}
// }

// type AllInbox struct {
// 	Map   map[string]*Inbox
// 	Mutex sync.Mutex
// }

// func NewAllInbox() *AllInbox {
// 	return &AllInbox{
// 		Map:   map[string]*Inbox{},
// 		Mutex: sync.Mutex{},
// 	}
// }
// func (allInbox *AllInbox) Get(look_for string) (*Inbox, bool) {
// 	allInbox.Mutex.Lock()
// 	defer allInbox.Mutex.Unlock()
// 	inbox, ok := allInbox.Map[look_for]
// 	return inbox, ok
// }

// func (allInbox *AllInbox) Add(key string, gs *Inbox) {
// 	allInbox.Mutex.Lock()
// 	defer allInbox.Mutex.Unlock()

// 	allInbox.Map[key] = gs
// }

// func (allInbox *AllInbox) Delete(look_for string) bool {
// 	allInbox.Mutex.Lock()
// 	defer allInbox.Mutex.Unlock()
// 	_, ok := allInbox.Map[look_for]
// 	delete(allInbox.Map, look_for)
// 	return ok
// }

// func (app *APIApplication) inboxWebsocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
// 		OriginPatterns: CORS_PATTERNS,
// 	})
// 	if err != nil {
// 		slog.Error("error accepting new connection", "err", err)
// 		return
// 	}
// 	defer conn.Close(websocket.StatusNormalClosure, "websocket ended")

// 	socketClient := socket.NewSocketClient(conn)

// 	personRow, err := ContextGetAuthenticatedUser(app.Queries, r)
// 	if err != nil {
// 		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_AUTH)
// 		return
// 	}

// 	socketClient.ListenInBackground(r)
// 	inbox := NewInbox(socketClient)
// 	app.InboxMap.Add(strconv.Itoa(personRow.ID), inbox)
// 	defer app.InboxMap.Delete(strconv.Itoa(personRow.ID))

// 	// socketSubber := rooms.NewRoomSocketListener(roomPlayer)
// 	// currentRoom.Subscribe(socketSubber)
// 	// defer currentRoom.Unsubscribe(socketSubber)

// 	for {
// 		select {
// 		case <-socketClient.Ctx.Done():
// 			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
// 			return
// 		case newMsgBytes := <-inbox.SocketClient.EventChan:
// 			newRoomMsg := &chatv1.InboxRequests{}
// 			err := proto.Unmarshal(newMsgBytes, newRoomMsg)
// 			if err != nil {
// 				slog.Error("can't unmarshal room msg", "error", err)
// 				continue
// 			}
// 			switch newRoomMsg.Content.(type) {
// 			case *chatv1.InboxRequests_SendMessageRequest:
// 				app.handleSendMessage(newRoomMsg.GetSendMessageRequest(), personRow.ID)
// 			default:
// 				slog.Error("unknown message type", "type", newRoomMsg.Content)
// 				continue
// 			}
// 		}
// 	}
// }

// func (app *APIApplication) handleSendMessage(sendMessageRequest *chatv1.SendMessageRequest, senderPersonID int) error {
// 	content := sendMessageRequest.GetContent()
// 	if content == "" {
// 		slog.Error("content is empty")
// 		return
// 	}

// 	chatID := int(req.Msg.ChatId)

// 	ctx := context.Background()

// 	// Verify user is a participant in the chat
// 	participant, err := app.Queries.IsChatParticipant(ctx, database.IsChatParticipantParams{
// 		ChatRoomIDRef: chatID,
// 		PersonID:      senderPersonID,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Ensure valid participant row found
// 	if participant.ParticipantID == 0 {
// 		return nil, connect.NewError(connect.CodeNotFound, errors.New("you are not a participant of this chat"))
// 	}

// 	// Insert the new message
// 	messageRow, err := app.Queries.InsertChatMessage(ctx, database.InsertChatMessageParams{
// 		ChatRoomIDRef:  chatID,
// 		SenderPersonID: senderPersonID,
// 		Content:        req.Msg.Content,
// 	})
// 	if err != nil {
// 		return nil, connect.NewError(connect.CodeInternal, err)
// 	}

// }
