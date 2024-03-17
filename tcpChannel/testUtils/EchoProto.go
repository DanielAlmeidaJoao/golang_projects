package testUtils

import (
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	"log"
	"net"
)

type ProtoEcho struct {
	id               gobabelUtils.APP_PROTO_ID
	ChannelInterface protoListener.ChannelInterface
	ServerAddr       string
}

func (p *ProtoEcho) ProtocolUniqueId() gobabelUtils.APP_PROTO_ID {
	return p.id
}

func (p *ProtoEcho) OnStart(channelInterface protoListener.ChannelInterface) {
	p.ChannelInterface = channelInterface
	log.Println("PROTOCOL STARTED: ", p.ProtocolUniqueId())
}
func (p *ProtoEcho) OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface protoListener.ChannelInterface) {
	str := string(msg)
	fmt.Println("------------ RECEIVED MESSAGE FROM: -----------------", (*from).String(), str)
}
func (p *ProtoEcho) ConnectionUp(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Printf("CONNECTION IS UP. FROM <%s>\n", (*from).String())
	p.ServerAddr = (*from).String()
}
func (p *ProtoEcho) ConnectionDown(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Printf("CONNECTION IS DOW. FROM <%s>\n", (*from).String())
}
func (p *ProtoEcho) HandleMessage(from string, protoSource gobabelUtils.APP_PROTO_ID, data protoListener.NetworkMessage) {
	fmt.Println("RECEIVED A MESSAGE FROM ", from)
}

func (p *ProtoEcho) SendMessage(address string, data string) (int, error) {
	return p.ChannelInterface.SendAppData(address, p.ProtocolUniqueId(), p.ProtocolUniqueId(), []byte(data))
}
func NewEchoProto() *ProtoEcho {
	return &ProtoEcho{
		id: gobabelUtils.APP_PROTO_ID(45),
	}
}
