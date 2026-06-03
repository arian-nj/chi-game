package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"connectrpc.com/connect"
	authv1 "github.com/arian-nj/chigame/backend/gen/auth/v1"
)

func (app *APIApplication) ValidateGuest(
	ctx context.Context,
	req *connect.Request[authv1.ValidateGuestRequest],
) (*connect.Response[authv1.ValidateGuestResponse], error) {
	deviceID := req.Msg.GetDeviceId()
	if deviceID == "" {
		deviceID = generateSecureDeviceID()
	}

	person, err := app.GetOrCreateGuestUser(ctx, deviceID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	tokenString, err := createGuestToken(person.ID, deviceID, app.Config.Jwt.SecretKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("sign token: %w", err))
	}

	return connect.NewResponse(&authv1.ValidateGuestResponse{
		DeviceId: deviceID,
		Token:    tokenString,
	}), nil
}

func generateSecureDeviceID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
