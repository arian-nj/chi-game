package api

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	chatv1 "github.com/arian-nj/chigame/backend/gen/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetAllChats returns all chats of the user
func (app *APIApplication) GetAllChats(ctx context.Context, req *connect.Request[chatv1.GetAllChatsRequest]) (*connect.Response[chatv1.GetAllChatsResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	chatsRows, err := app.Queries.GetAllChatsOfUser(ctx, personRow.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	chats := []*chatv1.Chat{}
	for _, chat := range chatsRows {
		chats = append(chats, &chatv1.Chat{
			Id:        int64(chat.ID),
			ChatType:  string(chat.Type),
			Name:      chat.Name.String,
			UpdatedAt: timestamppb.New(chat.UpdatedAt.Time),
		})
	}

	return connect.NewResponse(&chatv1.GetAllChatsResponse{
		Chats: []*chatv1.Chat{},
	}), nil
}

// GetMessages returns all messages of a chat
func (app *APIApplication) GetMessages(ctx context.Context, req *connect.Request[chatv1.GetMessagesRequest]) (*connect.Response[chatv1.GetMessagesResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	// check if the user is a participant of the chat
	participantRow, err := app.Queries.IsChatParticipant(ctx, database.IsChatParticipantParams{
		ChatID: int(req.Msg.ChatId),
		UserID: int(personRow.ID),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if participantRow.ID == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("you are not a participant of this chat"))
	}

	// get messages of the chat
	messagesRows, err := app.Queries.GetMessagesOfChat(ctx, int(req.Msg.ChatId))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	messages := []*chatv1.Message{}
	for _, message := range messagesRows {
		messages = append(messages, &chatv1.Message{
			Id:        int64(message.ID),
			ChatId:    int64(message.ChatID),
			Content:   message.Content,
			SentAt:    timestamppb.New(message.SentAt.Time),
			DeletedAt: timestamppb.New(message.DeletedAt.Time),
		})
	}

	return connect.NewResponse(&chatv1.GetMessagesResponse{
		Messages: messages,
	}), nil
}

func (app *APIApplication) createChat(ctx context.Context, userOneID int, userTwoID int, createdBy int) (*database.Chat, error) {
	// check if the chat already exist
	chatId, err := app.Queries.DoesChatExist(ctx, database.DoesChatExistParams{
		UserID:   userOneID,
		UserID_2: userTwoID,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if chatId != 0 {
		slog.Error("chat already exist failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("chat already exist failed"))
	}

	// check is allowed to create a chat (check friendship)
	isFriends, err := app.Queries.CheckFriendship(ctx, database.CheckFriendshipParams{
		UserID:   userOneID,
		FriendID: userTwoID,
	})
	if err != nil {
		slog.Error("can't check friendship", "error", err, "user_id", userOneID, "friend_id", userTwoID)
		return nil, connect.NewError(connect.CodeInternal, errors.New("can't check friendship"))
	}
	if !isFriends {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("you are not friends to create a chat"))
	}

	// create the chat
	chatRow, err := app.Queries.CreateChat(ctx, database.CreateChatParams{
		Type:      database.ChatTypeDirect,
		CreatedBy: int(createdBy),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// add chat participants
	userIDs := []int{userOneID, userTwoID}
	for _, userID := range userIDs {
		err = app.Queries.InsertChatParticipants(ctx, database.InsertChatParticipantsParams{
			ChatID: chatRow.ID,
			UserID: userID,
		})
		if err != nil {
			slog.Error("can't add chat participants", "error", err, "chat_id", chatRow.ID, "user_id", userID)
			return nil, connect.NewError(connect.CodeInternal, errors.New("can't add chat participants"))
		}
	}

	return &chatRow, nil
}

// SendMessage adds a message to a chat after validating authentication and chat participation
func (app *APIApplication) SendMessage(ctx context.Context, req *connect.Request[chatv1.SendMessageRequest]) (*connect.Response[chatv1.SendMessageResponse], error) {
	user := app.AuthenticateHeader(ctx, req.Header())
	if user == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	chatID := int(req.Msg.ChatId)
	senderID := int(user.ID)

	// Verify user is a participant in the chat
	participant, err := app.Queries.IsChatParticipant(ctx, database.IsChatParticipantParams{
		ChatID: chatID,
		UserID: senderID,
	})
	if err != nil {
		// If participant record not found, optionally try to create the chat, then return error
		if errors.Is(err, sql.ErrNoRows) {
			_, createErr := app.createChat(ctx, senderID, chatID, senderID)
			if createErr != nil {
				return nil, connect.NewError(connect.CodeInternal, createErr)
			}
			return nil, connect.NewError(connect.CodeNotFound, errors.New("chat not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Ensure valid participant row found
	if participant.ID == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("you are not a participant of this chat"))
	}

	// Insert the new message
	messageRow, err := app.Queries.InsertChatMessage(ctx, database.InsertChatMessageParams{
		ChatID:       chatID,
		SenderUserID: senderID,
		Content:      req.Msg.Content,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Build and return the response
	resp := &chatv1.SendMessageResponse{
		Message: &chatv1.Message{
			Id:        int64(messageRow.ID),
			ChatId:    int64(messageRow.ChatID),
			Content:   messageRow.Content,
			SentAt:    timestamppb.New(messageRow.SentAt.Time),
			DeletedAt: timestamppb.New(messageRow.DeletedAt.Time),
		},
	}
	return connect.NewResponse(resp), nil
}
