package cmd

import (
	"encoding/binary"
	"errors"
)

type CommandCode uint32
type CommandType uint8

// 命令类型
const (
	Request  CommandType = 1
	Response CommandType = 2
)

// RPC命令对象
type Command interface {
	GetCommandType() CommandType
	GetCommandCode() CommandCode
	GetCorrelationId() uint32
	GetCommandBody() []byte
}

type command struct {
	commandType   CommandType
	commandCode   CommandCode
	correlationId uint32
	commandBody   []byte
}

func (c *command) GetCommandType() CommandType {
	return c.commandType
}

func (c *command) GetCommandCode() CommandCode {
	return c.commandCode
}

func (c *command) GetCorrelationId() uint32 {
	return c.correlationId
}

func (c *command) GetCommandBody() []byte {
	return c.commandBody
}

func NewCommand(commandType CommandType, commandCode CommandCode,
	correlationId uint32, commandBody []byte) Command {
	return &command{
		commandType:   commandType,
		commandCode:   commandCode,
		correlationId: correlationId,
		commandBody:   commandBody,
	}
}

func DecodeCommand(buf []byte) (Command, error) {
	if len(buf) < 9 {
		return nil, errors.New("command len less than 9")
	}
	c := command{}
	c.commandType = CommandType(uint8(buf[0]))
	c.commandCode = CommandCode(binary.BigEndian.Uint32(buf[1:5]))
	c.correlationId = binary.BigEndian.Uint32(buf[5:9])
	if len(buf) > 9 {
		c.commandBody = buf[9:]
	} else {
		c.commandBody = []byte{}
	}
	return &c, nil
}

func EncodeCommand(command Command) []byte {
	buf := make([]byte, 0)
	buf = append(buf, byte(command.GetCommandType()))
	buf = binary.BigEndian.AppendUint32(buf, uint32(command.GetCommandCode()))
	buf = binary.BigEndian.AppendUint32(buf, command.GetCorrelationId())
	if len(command.GetCommandBody()) > 0 {
		buf = append(buf, command.GetCommandBody()...)
	}
	return buf
}
