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
func (p *ProtoEcho) HandleMessage(from string, protoSource gobabelUtils.APP_PROTO_ID, data *protoListener.CustomReader) {
	fmt.Println("RECEIVED A MESSAGE FROM ", from)
	msg := DeserializeData(data)
	println("RECEIVED MESSAGE IS:", msg.Data, msg.Count)
}
func (p *ProtoEcho) HandleMessage2(from string, protoSource gobabelUtils.APP_PROTO_ID, data *protoListener.CustomReader) {
	fmt.Println("RECEIVED A RANDOM MESSAGE FROM ", from)
	msg := DeserializeDataRandomMSG(data)
	println("RECEIVED RANDOM MESSAGE IS:", msg.Data, msg.Count)
	fmt.Println("RANDOM NUMBER IS :", msg.Count)
}

func (p *ProtoEcho) HandleLocalCommunication(sourceProto, destProto gobabelUtils.APP_PROTO_ID, data interface{}) {
	fmt.Println(p.id, " RECEIVED DATA ", sourceProto)
	fmt.Println(sourceProto, destProto, data)
}

func test() {
	msg := EchoMessage{}
	var netMsg protoListener.NetworkMessage = &msg
	msg2 := netMsg
	println(msg2)
}
func (p *ProtoEcho) SendMessage(address string, data string) (int, error) {
	return p.ChannelInterface.SendAppData(address, p.ProtocolUniqueId(), p.ProtocolUniqueId(), []byte(data))
}
func NewEchoProto() *ProtoEcho {
	return &ProtoEcho{
		id: gobabelUtils.APP_PROTO_ID(45),
	}
}
func NewEchoProto2() *ProtoEcho {
	return &ProtoEcho{
		id: gobabelUtils.APP_PROTO_ID(50),
	}
}
