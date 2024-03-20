package main

import (
	"encoding/binary"
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	testUtils "gobabel/testUtils"
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

/** ****/

// type MESSAGE_HANDLER_TYPE func(from string, protoSource APP_PROTO_ID, data protocolLIstenerLogics.NetworkMessage)
func main() {
	//go build ./...
	fmt.Println("SERVER STARTED")
	cc := make(chan int)
	pp := protoListener.NewProtocolListener("localhost", 3000, gobabelUtils.SERVER, binary.LittleEndian)
	proto := testUtils.NewEchoProto()
	fmt.Println("SERVER STARTED2")
	pp.AddProtocol(proto)
	fmt.Println("SERVER STARTED3")
	//(handlerId gobabelUtils.MessageHandlerID, funcHandler MESSAGE_HANDLER_TYPE, deserializer MESSAGE_DESERIALIZER_TYPE)
	pp.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(2), proto.HandleMessage)
	fmt.Println("SERVER STARTED4")
	pp.Start()
	fmt.Println("SERVER STARTED5")

	go func() {
		for {
			return
			time.Sleep(time.Second * 10)
			if proto.ChannelInterface.IsConnected(proto.ServerAddr) {
				now := time.Now()
				str := fmt.Sprintf("TIME HERE IS %s", now.String())
				result, err := proto.SendMessage(proto.ServerAddr, str)
				fmt.Println("CLIENT SENT THE MESSAGE. RESULT:", result, err, proto.ServerAddr)
			} else {
				fmt.Println("NOT CONNECTED. CAN NOT SEND MESSAGES...")
				return
			}
		}
	}()
	<-cc
}
