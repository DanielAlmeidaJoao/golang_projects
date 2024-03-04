package protocolLIstenerLogics

import (
	commons2 "fileServer/protocolLIstenerLogics/commons"
	"net"
)

type ConState int

// todo
const (
	IN_CONNECTION_UP ConState = iota
	OUT_CONNECTION_UP
	IN_CONNECTION_DOWN
	OUT_CONNECTION_DOWN
	MESSAGE_RECEIVED
)

type NetworkEvent struct {
	From      net.Addr
	Data      []byte
	protoId   int16
	eventType ConState
	eventId   commons2.EventID
}
type ConnectionState struct {
	address net.Addr
	state   ConState
	err     error
}
type ProtoListener struct {
	protocols map[int16]Protocol
}
