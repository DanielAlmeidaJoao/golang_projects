package main

import (
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	"log"
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
func (p *ProtoEcho) sendMessage(address string, data string) (int, error) {
	return p.channelInterface.SendAppData(address, p.ProtocolUniqueId(), p.ProtocolUniqueId(), []byte(data))
}
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
	log.Println("PROTOCOL STARTED: ", p.ProtocolUniqueId())
}
func (p *ProtoEcho) OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface protoListener.ChannelInterface) {
	str := string(msg)
	fmt.Println("------------ RECEIVED MESSAGE FROM: -----------------", (*from).String(), str)
}
func (p *ProtoEcho) ConnectionUp(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Printf("CONNECTION IS UP. FROM <%s>\n", (*from).String())
	p.serverAddr = (*from).String()
}
func (p *ProtoEcho) ConnectionDown(from *net.Addr, channelInterface protoListener.ChannelInterface) {
	fmt.Printf("CONNECTION IS DOW. FROM <%s>\n", (*from).String())
}

func main() {
	//go build ./...
	cc := make(chan int)
	pp := protoListener.NewProtocolListener("localhost", 3000, gobabelUtils.SERVER)
	proto := NewEchoProto()
	pp.AddProtocol(proto)
	pp.Start()

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
