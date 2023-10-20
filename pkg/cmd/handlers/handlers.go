package handlers

import (
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd/codes"
)

var (
	initFlag                                               = false                                        // 初始化标注
	commandHandlers map[cmd.CommandCode]cmd.CommandHandler = make(map[cmd.CommandCode]cmd.CommandHandler) //全局注册
)

func initCommandHandlers() {
	if initFlag {
		return
	}
	commandHandlers[codes.OpenRoom] = openRoomHandlerInstance
	initFlag = true
}

func GetCommandHandlers() map[cmd.CommandCode]cmd.CommandHandler {
	initCommandHandlers()
	return commandHandlers
}
