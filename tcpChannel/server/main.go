package main

import (
	"encoding/binary"
	"fmt"
	gobabelUtils "gobabel/commons"
	protoListener "gobabel/protocolLIstenerLogics"
	testUtils "gobabel/testUtils"
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
	proto := testUtils.NewEchoProto(pp)
	proto2 := testUtils.NewEchoProto2(pp)

	fmt.Println("SERVER STARTED2")
	err := pp.AddProtocol(proto)
	fmt.Println(err)
	err = pp.AddProtocol(proto2)
	fmt.Println(err)

	fmt.Println("SERVER STARTED3")
	//(handlerId gobabelUtils.MessageHandlerID, funcHandler MESSAGE_HANDLER_TYPE, deserializer MESSAGE_DESERIALIZER_TYPE)
	err1 := pp.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(2), proto.HandleMessage)
	err2 := pp.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(3), proto.HandleMessage2) //registar no server

	//err2 = pp.RegisterLocalCommunication(proto.ProtocolUniqueId(), proto2.ProtocolUniqueId(), 656, proto2.HandleLocalCommunication) //registar no server
	//fmt.Println(err2)
	fmt.Println("ERROR REGISTERING MSG HANDLERS:", err1, err2)
	fmt.Println("SERVER STARTED4")
	pp.Start()
	fmt.Println("SERVER STARTED5")

	<-cc
}
