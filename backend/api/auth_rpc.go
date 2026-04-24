package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"connectrpc.com/connect"
	authv1 "github.com/arian-nj/chigame/backend/gen/auth/v1"
)

func (app *ApiApplication) ValidateGuest(ctx context.Context,
	req *connect.Request[authv1.ValidateGuestRequest]) (*connect.Response[authv1.ValidateGuestResponse], error) {

	deviceID := req.Msg.DeviceId
	if deviceID == "" {
		deviceID = generateSecureDeviceID()
	}

	person, err := app.GetOrCreateGuestUser(deviceID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	token := createGuestToken(int(person.ID), deviceID)
	tokenString, err := token.SignedString(app.Config.Jwt.SecretKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &authv1.ValidateGuestResponse{
		DeviceId: deviceID,
		Token:    tokenString,
	}
	return connect.NewResponse(res), nil
}

func generateSecureDeviceID() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
