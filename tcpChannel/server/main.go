package main

import (
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	"net"
)

/**
type ProtoInterface interface {
	ProtocolUniqueId() gobabelUtils.APP_PROTO_ID
	OnStart(channelInterface *channel.ChannelInterface)
	OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface *channel.ChannelInterface)
	ConnectionUp(from *net.Addr, channelInterface *channel.ChannelInterface)
	ConnectionDown(from *net.Addr, channelInterface *channel.ChannelInterface)
}
*/

func NewEchoProto() *ProtoEcho {
	return &ProtoEcho{
		id: gobabelUtils.APP_PROTO_ID(45),
	}
}

type ProtoEcho struct {
	id               gobabelUtils.APP_PROTO_ID
	channelInterface protoListener.ChannelInterface
}

func (p *ProtoEcho) ProtocolUniqueId() gobabelUtils.APP_PROTO_ID {
	return p.id
}

func (p *ProtoEcho) OnStart(channelInterface protoListener.ChannelInterface) {
	p.channelInterface = channelInterface
}
func (p *ProtoEcho) OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface protoListener.ChannelInterface) {
	str := string(msg)
	fmt.Sprintf("RECEIVED MESSAGE FROM: ", (*from).String(), str)
}
func (p *ProtoEcho) ConnectionUp(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Sprintf("CONNECTION IS UP", (*from).String())
}
func (p *ProtoEcho) ConnectionDown(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Sprintf("CONNECTION IS DOW", (*from).String())
}

func main() {
	//go build ./...
	cc := make(chan int)
	pp := protoListener.NewProtocolListener("localhost", 3000, gobabelUtils.SERVER)
	proto := NewEchoProto()
	pp.AddProtocol(proto)
	pp.Start()

	<-cc
}
