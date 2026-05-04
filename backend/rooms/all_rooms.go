package rooms

import (
	"sync"
	"time"
)

type AllRooms struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

func NewAllRooms() *AllRooms {
	return &AllRooms{
		Rooms: map[string]*Room{},
		Mutex: sync.Mutex{},
	}
}

func (allRooms *AllRooms) Get(look_for string) (*Room, bool) {
	allRooms.Mutex.Lock()
	defer allRooms.Mutex.Unlock()
	room, ok := allRooms.Rooms[look_for]
	return room, ok
}

func (allRooms *AllRooms) Add(key string, gs *Room) {
	allRooms.Mutex.Lock()
	defer allRooms.Mutex.Unlock()

	allRooms.Rooms[key] = gs
}

func (allRooms *AllRooms) Delete(look_for string) bool {
	allRooms.Mutex.Lock()
	defer allRooms.Mutex.Unlock()
	_, ok := allRooms.Rooms[look_for]
	delete(allRooms.Rooms, look_for)
	return ok
}

func ClearDeadGamesCron(allRooms *AllRooms) {
	for {
		nowTime := time.Now()
		for key, gameRoom := range allRooms.Rooms {
			if nowTime.Sub(gameRoom.CreatedAt) > gameRoom.ExpireDuaration {
				allRooms.Mutex.Lock()
				delete(allRooms.Rooms, key)
				allRooms.Mutex.Unlock()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
