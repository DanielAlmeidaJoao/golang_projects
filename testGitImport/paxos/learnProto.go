package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
)

/****************************** LEARNER ****************************/

type LearnerProto struct {
	decided_value *PaxosMsg
	protoManager  tcpChannel.ProtoListenerInterface
	self          string
	currentTerm   uint32
	decided       map[uint32]*PaxosMsg
}

func NewLearnerProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *LearnerProto {
	return &LearnerProto{
		protoManager: listenerInterface,
		self:         address,
		currentTerm:  1,
		decided:      make(map[uint32]*PaxosMsg, 1000),
	}
}
func (receiver *LearnerProto) onDecided(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	value := ReadPaxosMsg(data)
	if receiver.currentTerm == value.term {
		for {
			receiver.currentTerm++
			//log.Println(receiver.self, " WWWWWWW DECISION TAKEN ", value)
			// func(sourceProto APP_PROTO_ID, destProto APP_PROTO_ID, data interface{})
			_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), CLIENT_PROTO_ID, value, ValueDecided)                //registar no server
			_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), ACCEPTOR_PROTO_ID, value.term, AcceptorValueDecided) //registar no server

			receiver.decided_value = value
			value = receiver.decided[receiver.currentTerm]
			if value == nil {
				return
			}
		}
		//TODO sendRequestToClient
	} else if receiver.currentTerm < value.term {
		receiver.decided[value.term] = value
	}
}

func (a *LearnerProto) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return LEARNER_PROTO_ID
}
func (a *LearnerProto) OnStart(channelInterface tcpChannel.ChannelInterface) {
	err1 := a.protoManager.RegisterNetworkMessageHandler(ON_DECIDE_ID, a.onDecided) //registar no server
	log.Println("REGISTERED ON DECIDED: ", err1)
}
func (a *LearnerProto) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *LearnerProto) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
func (a *LearnerProto) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
