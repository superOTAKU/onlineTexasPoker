package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd/handlers"
	"github.com/superOTAKU/onlineTexasPoker/pkg/config"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
	"github.com/superOTAKU/onlineTexasPoker/pkg/server"
)

var (
	cfgFile string
)

// cobra：帮助解析命令行参数的类库
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "online texas server manager",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			panic(err)
		}
		fmt.Printf("config loaded: %v\n", cfg)
		if err := log.InitLogger(&cfg.Log); err != nil {
			panic(err)
		}
		fmt.Printf("logger init\n")
		var protocolType server.ProtocolType
		switch cfg.Server.ProtocolType {
		case "Tcp":
			protocolType = server.Tcp
		case "WebSocket":
			protocolType = server.WebSocket
		default:
			panic("invalid protocol type : %v\n" + cfg.Server.ProtocolType)
		}
		serverOptions := &server.ServerOptions{
			ProtocolType:    protocolType,
			Host:            cfg.Server.Host,
			Port:            cfg.Server.Port,
			CommandHandlers: handlers.GetCommandHandlers(),
		}
		s := server.NewServer(serverOptions)
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file")
	rootCmd.Execute()
}
