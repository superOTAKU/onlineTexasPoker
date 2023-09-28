package server

// 封装基于TCP协议，定长报文的通讯协议细节，供上层调用

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log/logFields"
	"go.uber.org/zap"
)

type tcpProtocol struct {
	server      Server
	listener    net.Listener
	connMutex   sync.Mutex
	clientConns map[int32]*tcpConn
	clientMaxId atomic.Int32
}

var _ protocol = (*tcpProtocol)(nil)

func (p *tcpProtocol) getServer() Server {
	return p.server
}

type tcpConn struct {
	connId   int32
	conn     net.Conn
	protocol *tcpProtocol
}

var _ cmd.CommandContext = (*tcpConn)(nil)
var _ cmd.ClientConn = (*tcpConn)(nil)

func (c *tcpConn) GetConn() cmd.ClientConn {
	return c
}

func (c *tcpConn) WriteCommand(command cmd.Command) error {
	buf := cmd.EncodeCommand(command)
	packet := make([]byte, 0)
	packet = binary.BigEndian.AppendUint32(packet, uint32(len(buf)))
	packet = append(packet, buf...)
	if _, err := c.conn.Write(packet); err != nil {
		return err
	}
	return nil
}

// 定长字节编码，读一个完整数据包
func (c *tcpConn) readPacket() ([]byte, error) {
	packetLenBuf := make([]byte, 4)
	if _, err := io.ReadFull(c.conn, packetLenBuf); err != nil {
		return nil, err
	}
	packetLen := binary.BigEndian.Uint32(packetLenBuf)
	packet := make([]byte, packetLen)
	if _, err := io.ReadFull(c.conn, packet); err != nil {
		return nil, err
	}
	return packet, nil
}

func (c *tcpConn) serve() {
	defer func() {
		if r := recover(); r != nil {
			log.Logger().Error("conn serve panic", zap.Any("panic", r))
		}
		if err := c.conn.Close(); err != nil {
			log.Logger().Error("fail close conn", logFields.RemoteAddr(c.conn), zap.Error(err))
		}
		c.protocol.RemoteConn(c)
	}()
	log.Logger().Info("conn start serve", logFields.RemoteAddr(c.conn))
	for {
		packet, err := c.readPacket()
		if err != nil {
			log.Logger().Error("fail read packet", logFields.RemoteAddr(c.conn), zap.Error(err))
			break
		}
		command, err := cmd.DecodeCommand(packet)
		if err != nil {
			log.Logger().Error("fail decode command", logFields.RemoteAddr(c.conn), zap.Error(err))
			break
		}
		handler := c.protocol.getServer().GetCommandHandlers()[command.GetCommandCode()]
		if handler == nil {
			log.Logger().Error("command handler not found",
				logFields.RemoteAddr(c.conn), zap.Any("commandCode", command.GetCommandCode()))
		}
		c.handleCommand(handler, command)
	}
}

func (c *tcpConn) handleCommand(handler cmd.CommandHandler, command cmd.Command) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger().Error("handle command panic",
				zap.Any("remoteAddr", c.conn.RemoteAddr()), zap.Any("commandCode", command.GetCommandCode()))
		}
	}()
	handler.Handle(c, command)
}

func (p *tcpProtocol) ListenAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.server.GetListenHost(), p.server.GetListenPort()))
	p.listener = listener
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Logger().Error("fail listen conn", zap.Error(err))
			continue
		}
		connWrapper := &tcpConn{
			connId:   p.clientMaxId.Add(1),
			conn:     conn,
			protocol: p,
		}
		log.Logger().Info("received new conn", logFields.RemoteAddr(conn))
		p.connMutex.Lock()
		p.clientConns[connWrapper.connId] = connWrapper
		p.connMutex.Unlock()
		go connWrapper.serve()
	}
}

func (p *tcpProtocol) RemoteConn(conn *tcpConn) {
	p.connMutex.Lock()
	defer p.connMutex.Unlock()
	p.clientConns[conn.connId] = nil
}

func newTcpProtocol(server Server) protocol {
	p := &tcpProtocol{
		server:      server,
		connMutex:   sync.Mutex{},
		clientConns: make(map[int32]*tcpConn),
		clientMaxId: atomic.Int32{},
	}
	return p
}
