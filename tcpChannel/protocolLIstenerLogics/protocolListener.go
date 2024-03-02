package protocolLIstenerLogics

import (
	"fileServer/commons"
	"net"
)

type CON_STATE int

const (
	IN_CONNECTION_UP CON_STATE = iota
	OUT_CONNECTION_UP
	IN_CONNECTION_DOWN
	OUT_CONNECTION_DOWN
)

type ConnectionState struct {
	address net.Addr
	state   CON_STATE
	err     error
}
type ProtoListener struct {
	conStateChannel  chan ConnectionState // will receive the state of a connection, connection up or connection down
	inMessageChannel chan commons.DataWrapper
}
