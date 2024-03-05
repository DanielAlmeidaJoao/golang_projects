package protocolLIstenerLogics

import (
	"net"
)

type ConState int
type eventHandlerId int //protocol handler id

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
	ProtoId   int16
	EventType ConState
	EventId   eventHandlerId
}

func NewNetworkEvent(from net.Addr, data []byte, protoId int16, eventType ConState, eventId eventHandlerId) *NetworkEvent {
	return &NetworkEvent{
		From:      from,
		Data:      data,
		ProtoId:   protoId,
		EventType: eventType,
		EventId:   eventId,
	}
}

type ConnectionState struct {
	address net.Addr
	state   ConState
	err     error
}
type ProtoListener struct {
	protocols map[int16]*Protocol
}

func (l *ProtoListener) DeliverEvent(event *NetworkEvent) {
	protocol := l.protocols[event.ProtoId]
	if protocol == nil {
		//TODO log errror msg. MSG NOT BEING PROCESSED
	} else {
		protocol.eventQueue <- event
	}
}
