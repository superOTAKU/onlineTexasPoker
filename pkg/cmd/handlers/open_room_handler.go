package handlers

import (
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd/codes"
	"github.com/superOTAKU/onlineTexasPoker/pkg/game"
)

var openRoomHandlerInstance = &openRoomHandler{}

type openRoomHandler struct{}

type OpenRoomResponse struct {
	RoomId string
}

func (h *openRoomHandler) Handle(context cmd.CommandContext, command cmd.Command) error {
	newRoomId := game.GetRoomManager().NewRoom()
	return context.GetConn().WriteCommand(
		cmd.NewCommand(
			cmd.Response,
			codes.OpenRoom,
			command.GetCorrelationId(),
			[]byte(newRoomId)))
}
