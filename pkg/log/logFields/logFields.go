package logFields

import (
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HasRemoteAddr interface {
	RemoteAddr() net.Addr
}

func RemoteAddr(h HasRemoteAddr) zapcore.Field {
	return zap.Any("remoteAddr", h.RemoteAddr())
}
