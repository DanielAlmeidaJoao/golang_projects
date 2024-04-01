package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
)

/****************************** LEARNER ****************************/

type LearnerProto struct {
	decided_value *PaxosMsg
	protoManager  tcpChannel.ProtoListenerInterface
}

func NewLearnerProtocol(listenerInterface tcpChannel.ProtoListenerInterface) *LearnerProto {
	return &LearnerProto{
		protoManager: listenerInterface,
	}
}
func (receiver *LearnerProto) onDecided(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	value := ReadPaxosMsg(data)
	log.Println("DECISION TAKEN ", value)
	// func(sourceProto APP_PROTO_ID, destProto APP_PROTO_ID, data interface{})
	_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), CLIENT_PROTO_ID, value, ValueDecided) //registar no server
	receiver.decided_value = value
	//TODO sendRequestToClient
}

func (a *LearnerProto) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return LEARNER_PROTO_ID
}
func (a *LearnerProto) OnStart(channelInterface tcpChannel.ChannelInterface) {
	log.Println("LEARNER STARTED")
	err1 := a.protoManager.RegisterNetworkMessageHandler(ON_DECIDE_ID, a.onDecided) //registar no server
	log.Println("REGISTERED ON DECIDED: ", err1)
}
func (a *LearnerProto) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *LearnerProto) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
func (a *LearnerProto) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
