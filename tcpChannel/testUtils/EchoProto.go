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
	Counter          int
	listener         protoListener.ProtoListenerInterface
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
	msg := DeserializeData(data)
	log.Println("RECEIVED MESSAGE IS:", msg.Data, msg.Count)
	pair := &protoListener.CustomPair[string, *EchoMessage]{
		First:  from,
		Second: msg,
	}
	err2 := p.listener.RegisterLocalCommunication(p.ProtocolUniqueId(), 50, pair, p.Proto2GoingToReply) //registar no server
	if err2 != nil {
		log.Println("ERROR SENDING DATA TO ANOTHER PROTO", err2)
	}
}
func (p *ProtoEcho) ClientHandleMessage(from string, protoSource gobabelUtils.APP_PROTO_ID, data *protoListener.CustomReader) {
	msg := DeserializeData(data)
	log.Println("CLIENT RECEIVED MESSAGE IS:", msg.Data, msg.Count, from)
}
func (p *ProtoEcho) HandleMessage2(from string, protoSource gobabelUtils.APP_PROTO_ID, data *protoListener.CustomReader) {
	fmt.Println("RECEIVED A RANDOM MESSAGE FROM ", from)
	msg := DeserializeDataRandomMSG(data)
	println("RECEIVED RANDOM MESSAGE IS:", msg.Data, msg.Count)
	fmt.Println("RANDOM NUMBER IS :", msg.Count)
}

func (p *ProtoEcho) HandleLocalCommunication(sourceProto, destProto gobabelUtils.APP_PROTO_ID, data interface{}) {
	log.Println(sourceProto, destProto, data)
}

func (p *ProtoEcho) Proto2GoingToReply(sourceProto, destProto gobabelUtils.APP_PROTO_ID, data interface{}) {
	log.Println(p.id, "RECEVEID DATA FROM:", sourceProto, destProto, data)

	pair, ok := data.(*protoListener.CustomPair[string, *EchoMessage])

	if ok {
		result, er := p.ChannelInterface.SendAppData2(pair.First, 45, 45, pair.Second, 2)
		log.Println("SERVER REPLIED WITH:", result, er)
	}
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
func NewEchoProto(manager protoListener.ProtoListenerInterface) *ProtoEcho {
	return &ProtoEcho{
		id:       gobabelUtils.APP_PROTO_ID(45),
		listener: manager,
	}
}
func NewEchoProto2(manager protoListener.ProtoListenerInterface) *ProtoEcho {
	return &ProtoEcho{
		id:       gobabelUtils.APP_PROTO_ID(50),
		listener: manager,
	}
}
