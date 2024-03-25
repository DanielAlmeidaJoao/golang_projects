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

func TimerFunc1(proto gobabelUtils.APP_PROTO_ID, message interface{}) {
	fmt.Println("TIMER TRIGERED!!!! !!!! !!!!")
	msg, ok := message.(*testUtils.EchoMessage)
	if ok {
		fmt.Println("message is: %s %d", msg.Data, msg.Count)
	} else {
		fmt.Println("DATA IS WRONG TYPE")
	}
}
func main() {
	// go build ./...
	fmt.Println("SERVER STARTED")

	cc := make(chan int)

	protocolsManager := protoListener.NewProtocolListener("localhost", 3002, gobabelUtils.SERVER, binary.LittleEndian)
	proto := testUtils.NewEchoProto()
	protocolsManager.AddProtocol(proto)
	err1 := protocolsManager.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(2), proto.HandleMessage)
	err2 := protocolsManager.RegisterNetworkMessageHandler(gobabelUtils.MessageHandlerID(3), proto.HandleMessage2) //registar no server
	fmt.Println("ERROR REGISTERING MSG HANDLERS:", err1, err2)

	protocolsManager.Start()

	time.Sleep(time.Second * 3)
	proto.ChannelInterface.OpenConnection("localhost", 3000, 45)
	time.Sleep(time.Second * 5)

	msg := testUtils.EchoMessage{
		Data:  "SAY WHAAAAT ??",
		Count: 134,
	}
	//SendAppData2(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error)
	result, er := proto.ChannelInterface.SendAppData2(proto.ServerAddr, 45, 45, &msg, 2)
	fmt.Println("RESULT IS: ", result, er)

	protocolsManager.RegisterTimeout(proto.ProtocolUniqueId(), time.Second*5, &msg, TimerFunc1)

	go func() {
		return
		count := 0
		alternate := false
		for {
			alternate = !alternate
			time.Sleep(time.Second * 5)
			var result int
			var err error
			if proto.ChannelInterface.IsConnected(proto.ServerAddr) {
				now := time.Now()
				str := fmt.Sprintf("TIME HERE IS %s", now.String())
				if alternate {
					msg = testUtils.EchoMessage{
						Data:  str,
						Count: int32(count),
					}
					result, err = proto.ChannelInterface.SendAppData2(proto.ServerAddr, 45, 45, &msg, 2)
				} else {
					cc := float32(count) * 1.5
					msg2 := testUtils.RandomMessage{
						Data:  str,
						Count: cc,
					}
					result, err = proto.ChannelInterface.SendAppData2(proto.ServerAddr, 45, 45, &msg2, 3)
				}

				count++
				fmt.Println("CLIENT SENT THE MESSAGE. RESULT:", result, err, proto.ServerAddr)
			} else {
				fmt.Println("NOT CONNECTED. CAN NOT SEND MESSAGES...")
				return
			}
		}
	}()
	<-cc
}
