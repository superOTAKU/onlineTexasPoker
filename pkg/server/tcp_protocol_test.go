package server

import (
	"encoding/binary"
	"io"
	"net"
	"testing"

	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/config"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
)

// 测试用命令
type testCommandHandler struct{}

func (c *testCommandHandler) Handle(context cmd.CommandContext, command cmd.Command) error {
	context.GetConn().WriteCommand(cmd.NewCommand(cmd.Response, cmd.CommandCode(1),
		command.GetCorrelationId(), []byte("Hello")))
	return nil
}

var testCommandHandlers = make(map[cmd.CommandCode]cmd.CommandHandler)

// 测试TCP协议，发送请求并接收响应
func TestTcpProtocol(t *testing.T) {
	// 初始化日志到控制台
	log.InitLogger(&config.LogConfig{
		Console:  true,
		FilePath: "",
	})
	// 实现测试用CommandHandler
	testCommandHandlers[cmd.CommandCode(1)] = &testCommandHandler{}
	s := NewServer(&ServerOptions{
		ProtocolType:    Tcp,
		Host:            "localhost",
		Port:            1000,
		CommandHandlers: testCommandHandlers,
	})
	// 监听请求
	go s.ListenAndServe()
	// 发送请求，并等待响应
	conn, err := net.Dial("tcp", "localhost:1000")
	if err != nil {
		t.Fatalf("fail connect server: %v\n", err)
	} else {
		t.Logf("connect to server")
	}
	if n, err := conn.Write(binary.BigEndian.AppendUint32(make([]byte, 0), 14)); err != nil {
		t.Fatalf("fail write request len: %v\n", err)
	} else {
		t.Logf("write request len, n: %v", n)
	}
	if n, err := conn.Write(cmd.EncodeCommand(cmd.NewCommand(
		cmd.Request,
		cmd.CommandCode(1),
		1,
		[]byte("Hello")))); err != nil {
		t.Fatalf("fail write request body: %v\n", err)
	} else {
		t.Logf("write request body, n: %v", n)
	}
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		t.Fatalf("fail read response len: %v\n", err)
	}
	len := binary.BigEndian.Uint32(lenBuf)
	body := make([]byte, len)
	if _, err := io.ReadFull(conn, body); err != nil {
		t.Fatalf("fail read response body: %v\n", err)
	}
	command, err := cmd.DecodeCommand(body)
	if err != nil {
		t.Fatalf("decode command body error: %v\n", err)
	}
	t.Logf("response command is: %v\n", command)
}
