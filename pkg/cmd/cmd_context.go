package cmd

// 命令上下文
type CommandContext interface {
	// 获取当前连接
	GetConn() ClientConn
}

// 连接层抽象
type ClientConn interface {
	// 回写命令
	WriteCommand(command Command) error
}
