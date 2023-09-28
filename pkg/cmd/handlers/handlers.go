package handlers

import (
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
)

var (
	initFlag                                               = false                                        // 初始化标注
	commandHandlers map[cmd.CommandCode]cmd.CommandHandler = make(map[cmd.CommandCode]cmd.CommandHandler) //全局注册
)

func initCommandHandlers() {
	if initFlag {
		return
	}
	initFlag = true
}

func GetCommandHandlers() map[cmd.CommandCode]cmd.CommandHandler {
	initCommandHandlers()
	return commandHandlers
}
