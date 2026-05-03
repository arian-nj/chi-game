package gamesessions

import (
	"strconv"
	"sync"
	"time"
)

type AllSession struct {
	Sessions map[string]*GameSession
	Mutex    sync.Mutex
}

func NewAllSessions() *AllSession {
	return &AllSession{
		Sessions: map[string]*GameSession{},
		Mutex:    sync.Mutex{},
	}
}

func (allSessions *AllSession) Get(look_for string) (*GameSession, bool) {
	allSessions.Mutex.Lock()
	defer allSessions.Mutex.Unlock()
	gameSession, ok := allSessions.Sessions[look_for]
	return gameSession, ok
}

func (allSessions *AllSession) IsSessionEmpty(playerId int) bool {
	allSessions.Mutex.Lock()
	defer allSessions.Mutex.Unlock()

	_, isFound := allSessions.Sessions[strconv.Itoa(playerId)]

	return !isFound
}

func (allSession *AllSession) Add(key string, gs *GameSession) {
	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	allSession.Sessions[key] = gs
}

func (allSessions *AllSession) Delete(look_for string) bool {
	allSessions.Mutex.Lock()
	defer allSessions.Mutex.Unlock()
	_, ok := allSessions.Sessions[look_for]
	delete(allSessions.Sessions, look_for)
	return ok
}

func ClearDeadGamesCron(allSessions *AllSession) {
	for {
		nowTime := time.Now()
		for key, gameSession := range allSessions.Sessions {
			if nowTime.Sub(gameSession.CreatedAt) > gameSession.ExpireDuaration {
				allSessions.Mutex.Lock()
				delete(allSessions.Sessions, key)
				allSessions.Mutex.Unlock()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
