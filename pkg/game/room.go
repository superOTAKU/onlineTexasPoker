package game

import "github.com/google/uuid"

var roomManager = &RoomManager{
	roomMap: make(map[string]*Room),
}

func GetRoomManager() *RoomManager {
	return roomManager
}

type RoomManager struct {
	roomMap map[string]*Room
}

type Room struct {
	id string
}

func (rm *RoomManager) NewRoom() string {
	id := uuid.NewString()
	rm.roomMap[id] = &Room{
		id: id,
	}
	return id
}
