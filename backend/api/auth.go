package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/golang-jwt/jwt/v5"
)

type AuthType string

const (
	AuthTypeGuestDevice = "guest_device"
)

func (app *ApiApplication) GetOrCreateGuestUser(deviceID string) (*database.Person, error) {
	person, err := app.Queries.GetPersonByAuthMethod(context.Background(), database.GetPersonByAuthMethodParams{
		AuthType:  AuthTypeGuestDevice,
		AuthValue: deviceID,
	})
	if err == nil {
		return &person, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query guest user: %w", err)
	}

	newPerson, err := app.Queries.InsertPerson(context.Background(), database.InsertPersonParams{
		IsGuest:  true,
		Username: random.GenerateRandomUsername(8),
	})
	if err != nil {
		return nil, fmt.Errorf("create person: %w", err)
	}

	// Link auth method
	_, err = app.Queries.InsertAuthMethod(context.Background(), database.InsertAuthMethodParams{
		UserID:    newPerson.ID,
		AuthType:  AuthTypeGuestDevice,
		AuthValue: deviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("link auth method: %w", err)
	}

	return &newPerson, nil

}

const JWTExpiryDuration = 1 * time.Hour
const GuestJWTExpiryDuration = 365 * 24 * time.Hour

type CustomClaim struct {
	IsGuest  bool   `json:"igu"`
	DeviceID string `json:"did,omitempty"`
	jwt.RegisteredClaims
}

func createRegisteredUserToken(userId int) *jwt.Token {
	claims := CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userId),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpiryDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		IsGuest: false,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func createGuestToken(userId int, deviceID string) *jwt.Token {
	claims := CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userId),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(GuestJWTExpiryDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		IsGuest:  true,
		DeviceID: deviceID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (app *ApiApplication) ValidateToken(tokenString string) (int, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (any, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("error parsing token", "err", err)
		return 0, err
	}

	claims, ok := token.Claims.(*CustomClaim)
	if !ok {
		slog.Error("error converting Claims to CustomClaims")
		return 0, err
	}

	expireAt, err := claims.GetExpirationTime()
	if err != nil {
		slog.Error("error getting expiration time", "err", err)
		return 0, err
	}
	if expireAt.Time.Unix() < time.Now().Unix() {
		slog.Error("token expired", "expireAt", expireAt)
		return 0, err
	}

	notBefore, err := claims.GetNotBefore()
	if err != nil {
		slog.Error("error getting not before", "err", err)
		return 0, err
	}

	if notBefore.Time.Unix() > time.Now().Unix() {
		slog.Error("token not before is in the future", "notBefore", notBefore)
		return 0, err
	}

	sub, err := claims.GetSubject()
	if err != nil {
		slog.Error("error getting subject", "err", err)
		return 0, err
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		slog.Error("error converting subject to int", "err", err)
		return 0, err
	}
	// if user.ID == 0 {
	// 	app.InvalidAuthenticationToken(w, r)
	// }

	return userID, nil
}

// type WebAppUser struct {
// 	TgID                  int64  `json:"id"`
// 	IsBot                 bool   `json:"is_bot"`
// 	FirstName             string `json:"first_name"`
// 	LastName              string `json:"last_name"`
// 	Username              string `json:"username"`
// 	LanguageCode          string `json:"language_code"`
// 	IsPremium             bool   `json:"is_premium"`
// 	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
// 	AllowsWriteToPM       bool   `json:"allows_write_to_pm"`
// 	PhotoURL              string `json:"photo_url"`
// }

// func ValidateWebappRequest(values url.Values, token string) (user *WebAppUser, ok bool) {
// 	h := values.Get("hash")
// 	values.Del("hash")

// 	var vals []string

// 	var u WebAppUser

// 	for k, v := range values {
// 		vv, _ := url.QueryUnescape(v[0])
// 		vals = append(vals, k+"="+vv)
// 		if k == "user" {
// 			errDecodeUser := json.Unmarshal([]byte(vv), &u)
// 			if errDecodeUser != nil {
// 				return nil, false
// 			}
// 		}
// 	}

// 	sort.Slice(vals, func(i, j int) bool {
// 		return vals[i] < vals[j]
// 	})

// 	hmac1 := hmac.New(sha256.New, []byte("WebAppData"))
// 	hmac1.Write([]byte(token))
// 	r1 := hmac1.Sum(nil)

// 	data := []byte(strings.Join(vals, "\n"))

// 	hmac2 := hmac.New(sha256.New, r1)
// 	hmac2.Write(data)
// 	r2 := hmac2.Sum(nil)

// 	if h != fmt.Sprintf("%x", r2) {
// 		return nil, false
// 	}

// 	return &u, true
// }
