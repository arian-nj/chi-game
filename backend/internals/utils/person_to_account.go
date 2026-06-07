package utils

import (
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
)

func PersonToAccount(person *database.Person) *accountv1.Account {
	return &accountv1.Account{
		Id:          person.ID,
		Username:    person.Username,
		DisplayName: person.DisplayName,
	}
}
