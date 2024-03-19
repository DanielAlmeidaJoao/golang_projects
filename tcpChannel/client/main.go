package main

import (
	"encoding/binary"
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	testUtils "gobabel/testUtils"
	"time"
)

//go build ./...
/**
type ProtoInterface interface {
	ProtocolUniqueId() gobabelUtils.APP_PROTO_ID
	OnStart(channelInterface *channel.ChannelInterface)
	OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface *channel.ChannelInterface)
	ConnectionUp(from *net.Addr, channelInterface *channel.ChannelInterface)
	ConnectionDown(from *net.Addr, channelInterface *channel.ChannelInterface)
}
*/

func main() {
	// go build ./...
	fmt.Println("SERVER STARTED")

	cc := make(chan int)

	protocolsManager := protoListener.NewProtocolListener("localhost", 3002, gobabelUtils.SERVER, binary.LittleEndian)
	proto := testUtils.NewEchoProto()
	protocolsManager.AddProtocol(proto)
	protocolsManager.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(2), proto.HandleMessage, testUtils.DeserializeData)
	protocolsManager.Start()

	time.Sleep(time.Second * 3)
	proto.ChannelInterface.OpenConnection("localhost", 3000, 45)
	time.Sleep(time.Second * 5)

	_ = testUtils.EchoMessage{
		Data:  "SAY WHAAAAT ??",
		Count: 134,
	}
	//	SendAppData2(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error)
	//result, er := proto.ChannelInterface.SendAppData2(proto.ServerAddr, 45, 45, &msg, 2)
	//fmt.Println("RESULT IS: ", result, er)
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
