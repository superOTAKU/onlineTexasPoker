package cmd

// 命令处理接口
type CommandHandler interface {
	Handle(context CommandContext, command Command) error
}
