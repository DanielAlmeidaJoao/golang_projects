package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"time"
)

/*
ProtocolUniqueId() APP_PROTO_ID
	OnStart(channelInterface ChannelInterface)
	OnMessageArrival(customCon *CustomConnection, source, destProto APP_PROTO_ID, msg []byte, channelInterface ChannelInterface)
	ConnectionUp(customCon *CustomConnection, channelInterface ChannelInterface)
	ConnectionDown(customCon *CustomConnection, channelInterface ChannelInterface)

*/

type ClientProtocol struct {
	ops          map[string]*PaxosMsg
	count        int
	protoManager tcpChannel.ProtoListenerInterface
	proper       *ProposerProtocol
}

func NewClientProtocol(listenerInterface tcpChannel.ProtoListenerInterface, proposer *ProposerProtocol) *ClientProtocol {
	return &ClientProtocol{
		ops:          make(map[string]*PaxosMsg),
		protoManager: listenerInterface,
		count:        0,
		proper:       proposer,
	}
}
func (c *ClientProtocol) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return CLIENT_PROTO_ID
}

func (c *ClientProtocol) handleTimer(sourceProto tcpChannel.APP_PROTO_ID, data interface{}) {
	msg := &PaxosMsg{
		msgId:    time.Now().String(),
		msgValue: "daniel joao",
	}
	f := c.proper.SendPaxosRequest
	err2 := c.protoManager.RegisterLocalCommunication(c.ProtocolUniqueId(), PROPOSER_PROTO_ID, msg, f) //registar no server
	log.Println("ERROR REGISTERING PROTO", err2)
}

func (c *ClientProtocol) OnStart(channelInterface tcpChannel.ChannelInterface) {
	c.protoManager.RegisterTimeout(c.ProtocolUniqueId(), time.Second*5, nil, c.handleTimer)

}
func (c *ClientProtocol) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (c *ClientProtocol) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
func (c *ClientProtocol) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
