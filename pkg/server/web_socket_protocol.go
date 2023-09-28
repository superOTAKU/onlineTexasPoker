package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log/logFields"
	"go.uber.org/zap"
)

type webSocketProtocol struct {
	ginCtx *gin.Context
}

type webSocketConn struct {
	conn *websocket.Conn
	p    *webSocketProtocol
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *webSocketProtocol) getServer() Server {
	return nil
}

func (c *webSocketConn) GetConn() cmd.ClientConn {
	return nil
}

func (c *webSocketConn) WriteCommand(command cmd.Command) error {

	return nil
}

func (p *webSocketProtocol) ListenAndServe() error {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Logger().Error("fail upgrade websocket", logFields.RemoteAddr(conn), zap.Error(err))
			return
		}
		log.Logger().Info("received webSocket conn", logFields.RemoteAddr(conn))
		connWrapper := &webSocketConn{
			conn: conn,
		}
		go connWrapper.serve()
	})
	log.Logger().Info("starting ws server", zap.String("Host",
		p.getServer().GetListenHost()), zap.Int("Port", p.getServer().GetListenPort()))
	return r.Run(fmt.Sprintf("%s:%d", p.getServer().GetListenHost(), p.getServer().GetListenPort()))
}

func (c *webSocketConn) serve() {
	defer func() {
		if r := recover(); r != nil {
			log.Logger().Error("conn serve panic", zap.Any("panic", r))
		}

	}()
	log.Logger().Info("conn start serve", logFields.RemoteAddr(c.conn))
	for {
		// webSocket协议本身已处理了消息长度，分片问题，只需要读整包
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Logger().Error("fail read msg", logFields.RemoteAddr(c.conn), zap.Error(err))
			break
		}
		if msgType != websocket.BinaryMessage {
			log.Logger().Error("receive text msg", logFields.RemoteAddr(c.conn))
			break
		}
		command, err := cmd.DecodeCommand(msg)
		if err != nil {
			log.Logger().Error("fail decode command", zap.Error(err))
		}
		handler := c.p.getServer().GetCommandHandlers()[command.GetCommandCode()]
		if handler == nil {
			log.Logger().Error("command handler not found", zap.Any("commandCode", command.GetCommandCode()))
			break
		}
		c.handleCommand(handler, command)
	}
}

func (c *webSocketConn) handleCommand(handler cmd.CommandHandler, command cmd.Command) {
	defer func() {
		if r := recover(); r != nil {
			
		}
	}()
	handler.Handle(c, command)
}
