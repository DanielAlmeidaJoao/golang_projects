package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"math"
	"time"
)

/****************************** LEARNER ****************************/

type LearnerProto struct {
	decided_value *PaxosMsg
	protoManager  tcpChannel.ProtoListenerInterface
	self          string
	currentTerm   uint32
	decided       map[uint32]*PaxosMsg
	peers         map[string]*tcpChannel.CustomConnection
	totalReceived int
}

func NewLearnerProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *LearnerProto {
	return &LearnerProto{
		protoManager:  listenerInterface,
		self:          address,
		currentTerm:   1,
		decided:       make(map[uint32]*PaxosMsg, 1000),
		peers:         make(map[string]*tcpChannel.CustomConnection),
		totalReceived: 0,
	}
}

func (p *LearnerProto) majority() int {
	return int(math.Ceil(float64((len(p.peers) + 1) / 2)))
}

func (receiver *LearnerProto) onDecided(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	receiver.totalReceived++
	value := ReadPaxosMsg(data)
	if value.term < receiver.currentTerm {
		return
	}
	//log.Println("0000000000000000000000000000000000 <SELF,TERM,PROP_NUMBER> ", receiver.self, value.term, value.proposalNum)
	aux := receiver.decided[value.term]
	if aux == nil {
		receiver.decided[value.term] = value
		aux = value
	}
	//log.Println("-------------LEARNED------------ <SELF,TERM,PROP_NUMBER,ID,COUNT--AUX_TERM> ", receiver.self, receiver.currentTerm, aux.proposalNum, aux.msgId, aux.decidedCount, aux.term)
	for aux != nil && aux.term == receiver.currentTerm {
		//log.Println("-------------LEARNED------------ <SELF,TERM,PROP_NUMBER,ID> ", receiver.self, receiver.currentTerm, aux.proposalNum, aux.msgId)
		delete(receiver.decided, aux.term)
		receiver.currentTerm++
		aux2 := aux
		_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), CLIENT_PROTO_ID, aux2, ValueDecided) //registar no server
		//_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), ACCEPTOR_PROTO_ID, aux2, AcceptorValueDecided) //DA ERRO, ACESSO CONCORRENTE DO MAPA
		//receiver.decided_value = value
		aux = receiver.decided[receiver.currentTerm]
	}
}

/*
	func (receiver *LearnerProto) onDecided(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
		receiver.totalReceived++
		value := ReadPaxosMsg(data)
		if value.term < receiver.currentTerm {
			return
		}
		//log.Println("0000000000000000000000000000000000 <SELF,TERM,PROP_NUMBER> ", receiver.self, value.term, value.proposalNum)
		aux := receiver.decided[value.term]
		if aux == nil || value.proposalNum > aux.proposalNum {
			receiver.decided[value.term] = value
			aux = value
		} else if value.proposalNum < aux.proposalNum {
			return
		}
		aux.decidedCount++
		log.Println("-------------LEARNED------------ <SELF,TERM,PROP_NUMBER,ID,COUNT--AUX_TERM> ", receiver.self, receiver.currentTerm, aux.proposalNum, aux.msgId, aux.decidedCount, aux.term)
		for aux != nil && aux.term == receiver.currentTerm && aux.decidedCount >= uint32(receiver.majority()) {
			//log.Println("-------------LEARNED------------ <SELF,TERM,PROP_NUMBER,ID> ", receiver.self, receiver.currentTerm, aux.proposalNum, aux.msgId)
			delete(receiver.decided, aux.term)
			receiver.currentTerm++
			aux2 := aux
			_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), CLIENT_PROTO_ID, aux2, ValueDecided) //registar no server
			//_ = receiver.protoManager.SendLocalEvent(receiver.ProtocolUniqueId(), ACCEPTOR_PROTO_ID, aux2, AcceptorValueDecided) //DA ERRO, ACESSO CONCORRENTE DO MAPA
			//receiver.decided_value = value
			aux = receiver.decided[receiver.currentTerm]
		}
	}
*/
func (a *LearnerProto) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return LEARNER_PROTO_ID
}
func (a *LearnerProto) OnStart(channelInterface tcpChannel.ChannelInterface) {
	err1 := a.protoManager.RegisterNetworkMessageHandler(ON_DECIDE_ID, a.onDecided) //registar no server
	log.Println("REGISTERED ON DECIDED: ", err1)
	a.protoManager.RegisterPeriodicTimeout(a.ProtocolUniqueId(), time.Second*30, nil, a.PeriodicTimerHandler)
}
func (a *LearnerProto) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *LearnerProto) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	a.peers[customCon.RemoteAddress().String()] = customCon
}
func (a *LearnerProto) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	log.Println("CONNECTION IS DOWN ", a.self)
}
func (c *LearnerProto) PeriodicTimerHandler(handlerId int, proto tcpChannel.APP_PROTO_ID, message interface{}) {
	log.Println("TIMER ON TRIGGERED LEARNER TOTAL RECEIVED", c.self, c.totalReceived)
}
