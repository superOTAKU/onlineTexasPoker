package server

import (
	"github.com/superOTAKU/onlineTexasPoker/pkg/cmd"
	"github.com/superOTAKU/onlineTexasPoker/pkg/log"
	"go.uber.org/zap"
)

type ProtocolType uint8

const (
	Tcp       ProtocolType = 1
	WebSocket ProtocolType = 2
)

type Server interface {
	ListenAndServe() error
	GetProtocol() ProtocolType
	GetListenHost() string
	GetListenPort() int
	GetCommandHandlers() map[cmd.CommandCode]cmd.CommandHandler
}

type protocol interface {
	getServer() Server
	ListenAndServe() error
}

type server struct {
	protocolType    ProtocolType
	host            string
	port            int
	commandHandlers map[cmd.CommandCode]cmd.CommandHandler
	protocol        protocol
}

func (s *server) ListenAndServe() error {
	log.Logger().Info("start listen", zap.String("host", s.host), zap.Int("port", s.port), zap.Any("protocolType", s.protocolType))
	if err := s.protocol.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *server) GetProtocol() ProtocolType {
	return s.protocolType
}

func (s *server) GetListenHost() string {
	return s.host
}

func (s *server) GetListenPort() int {
	return s.port
}

func (s *server) GetCommandHandlers() map[cmd.CommandCode]cmd.CommandHandler {
	return s.commandHandlers
}

type ServerOptions struct {
	ProtocolType    ProtocolType
	Host            string
	Port            int
	CommandHandlers map[cmd.CommandCode]cmd.CommandHandler
}

func NewServer(options *ServerOptions) Server {
	s := &server{
		protocolType:    options.ProtocolType,
		host:            options.Host,
		port:            options.Port,
		commandHandlers: options.CommandHandlers,
	}
	switch options.ProtocolType {
	case Tcp:
		s.protocol = newTcpProtocol(s)
	case WebSocket:
		s.protocol = newWebSocketProtocol(s)
	}
	return s
}
