package server

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/config"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
)

func TestWebSocketProtocol(t *testing.T) {
	// 初始化日志到控制台
	log.InitLogger(&config.LogConfig{
		Console:  true,
		FilePath: "",
	})
	// 实现测试用CommandHandler
	testCommandHandlers[cmd.CommandCode(1)] = &testCommandHandler{}
	s := NewServer(&ServerOptions{
		ProtocolType:    WebSocket,
		Host:            "localhost",
		Port:            1000,
		CommandHandlers: testCommandHandlers,
	})
	// 监听请求
	go s.ListenAndServe()
	// 等待gin起监听
	time.Sleep(time.Second * 3)
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:1000", nil)
	if err != nil {
		t.Fatalf("fail to connect ws server: %v", err)
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, cmd.EncodeCommand(cmd.NewCommand(
		cmd.Request,
		cmd.CommandCode(1),
		1,
		[]byte("Hello")))); err != nil {
		t.Fatalf("write request fail: %v", err)
	}
	_, bytes, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read message fail: %v", err)
	}
	command, err := cmd.DecodeCommand(bytes)
	if err != nil {
		t.Fatalf("decode command fail: %v", command)
	}
	t.Logf("received command: %v", command)
}
