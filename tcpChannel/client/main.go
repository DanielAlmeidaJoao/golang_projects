package main

import (
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	"net"
	"time"
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
	serverAddr       string
}

func (p *ProtoEcho) ProtocolUniqueId() gobabelUtils.APP_PROTO_ID {
	return p.id
}

func (p *ProtoEcho) OnStart(channelInterface protoListener.ChannelInterface) {
	p.channelInterface = channelInterface
	fmt.Println("PROTOCOL STARTED!!!")
	p.channelInterface.OpenConnection("localhost", 3000, p.ProtocolUniqueId())
}
func (p *ProtoEcho) OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface protoListener.ChannelInterface) {
	str := string(msg)
	fmt.Printf("****************************** RECEIVED MESSAGE FROM: %s \n", (*from).String(), str)
}
func (p *ProtoEcho) ConnectionUp(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Printf("CONNECTION IS UP", (*from).String())
	p.serverAddr = (*from).String()
	str := "A BRAVE NEW WORLD!"
	res, er := p.channelInterface.SendAppData((*from).String(), p.ProtocolUniqueId(), p.ProtocolUniqueId(), []byte(str))
	fmt.Println("MESSAGE SENT: ", res, er, len(str))
}
func (p *ProtoEcho) sendMessage(address string, data string) (int, error) {
	return p.channelInterface.SendAppData(address, p.ProtocolUniqueId(), p.ProtocolUniqueId(), []byte(data))
}
func (p *ProtoEcho) ConnectionDown(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Println("CONNECTION IS DOWN--", (*from).String())
}

func main() {
	// go build ./...
	cc := make(chan int)

	protocolsManager := protoListener.NewProtocolListener("localhost", 3001, gobabelUtils.SERVER)
	proto := NewEchoProto()
	protocolsManager.AddProtocol(proto)
	protocolsManager.Start()

	go func() {
		for {
			time.Sleep(time.Second * 10)
			if proto.channelInterface.IsConnected(proto.serverAddr) {
				now := time.Now()
				str := fmt.Sprintf("TIME HERE IS %s", now.String())
				result, err := proto.sendMessage(proto.serverAddr, str)
				fmt.Println("CLIENT SENT THE MESSAGE. RESULT:", result, err, proto.serverAddr)
			} else {
				fmt.Println("NOT CONNECTED. CAN NOT SEND MESSAGES...")
				return
			}
		}
	}()
	<-cc
}
