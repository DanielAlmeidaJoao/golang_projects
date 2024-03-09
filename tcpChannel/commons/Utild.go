package commons

import (
	"errors"
	"net"
)

var NOT_CONNECTED = errors.New("NOT_CONNECTED")
var PROTOCOL_EXIST_ALREADY = errors.New("PROTOCOL_EXIST_ALREADY")
var NO_PROTOCOLS_TO_RUN = errors.New("NO_PROTOCOLS_TO_RUN")

type MessageHandlerID int //protocol handler id
type APP_PROTO_ID uint16

const ALL_PROTO_ID = 1

// todo
type NET_EVENT int

const (
	CONNECTION_UP NET_EVENT = iota
	CONNECTION_DOWN
	MESSAGE_RECEIVED
)

type CONNECTION_TYPE int8

const (
	P2P CONNECTION_TYPE = iota
	CLIENT
	SERVER
)

type MSG_TYPE uint8

var UNEXPECTED_NETWORK_READ_DATA = errors.New("UNEXPECTED_NETWORK_READ_DATA")

const (
	LISTEN_ADDRESS_MSG MSG_TYPE = iota
	APP_MSG
)

type DataWrapper struct {
	From    net.Addr
	Data    []byte
	protoId int16
}

const MAX_BYTES_TO_READ = 2 * 1024

func Abort(err error) bool {
	if err == nil {
		return false
	} else {
		panic(err)
		return true
	}
}

type NetworkEvent struct {
	From             net.Addr
	Data             []byte
	SourceProto      APP_PROTO_ID
	DestProto        APP_PROTO_ID
	NET_EVENT        NET_EVENT
	MessageHandlerID MessageHandlerID
}

type MsgHandlerFunc func(from net.Addr, data []byte, sourceProto APP_PROTO_ID, eventType NET_EVENT)

func NewNetworkEvent(from net.Addr, data []byte, sourceProto, destProto APP_PROTO_ID, eventType NET_EVENT, messageHandlerID MessageHandlerID) *NetworkEvent {
	return &NetworkEvent{
		From:             from,
		Data:             data,
		SourceProto:      sourceProto,
		DestProto:        destProto,
		NET_EVENT:        eventType,
		MessageHandlerID: messageHandlerID,
	}
}
