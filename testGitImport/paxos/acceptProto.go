package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"time"
)

type termArg struct {
	accepted_num   uint32
	promised_num   uint32
	term           uint32
	accepted_value *PaxosMsg
	decided        *PaxosMsg
	remoteAddress  string
}

/************************************ ***************************************/
type AcceptorProto struct {
	terms        map[uint32]*termArg
	protoManager tcpChannel.ProtoListenerInterface
	peers        map[string]*tcpChannel.CustomConnection
	totalSent    int
	self         string
	proposer     string
}

func NewAcceptorProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *AcceptorProto {
	return &AcceptorProto{
		protoManager: listenerInterface,
		peers:        make(map[string]*tcpChannel.CustomConnection),
		totalSent:    0,
		self:         address,
		terms:        make(map[uint32]*termArg, 100000),
	}
}
func (a *AcceptorProto) onPrepare(selfProto tcpChannel.ProtoInterface, customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	preparaMessage := ReadDataPrepareMessage(data)
	arg := a.terms[preparaMessage.term]
	if arg == nil {
		arg = &termArg{
			term: preparaMessage.term,
		}
		a.terms[preparaMessage.term] = arg
	}
	if arg.promised_num < preparaMessage.proposal_num {
		//log.Println("2222222222222222222222222222 33333333333333333 ACCEPTING PREPARE: <SELF,TERM,PROSAL_NUM>", a.self, preparaMessage.term, preparaMessage.proposal_num)
		arg.promised_num = preparaMessage.proposal_num
		arg.remoteAddress = customConn.RemoteAddress().String()
		prom := &Promise{
			accepted_value: arg.accepted_value,
			accepted_num:   arg.accepted_num,
			promised_num:   preparaMessage.proposal_num,
			term:           arg.term,
		}
		customConn.SendData2(ACCEPTOR_PROTO_ID, PROPOSER_PROTO_ID, prom, ON_PROMISE_ID)
	}
}

func (a *AcceptorProto) setTermArg(accptMessage *AcceptMessage) *termArg {
	arg := a.terms[accptMessage.term]
	if arg == nil {
		arg = &termArg{
			term:           accptMessage.term,
			accepted_value: accptMessage.value,
		}
		a.terms[accptMessage.term] = arg
	}
	return arg
}

func (a *AcceptorProto) onAccept(protoInterface tcpChannel.ProtoInterface, customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	accptMessage := ReadDataAcceptMessage(data)
	arg := a.setTermArg(accptMessage)
	if arg.promised_num < accptMessage.proposal_num || (arg.promised_num == accptMessage.proposal_num && arg.remoteAddress <= customConn.RemoteAddress().String()) {
		//log.Println("PROMISE 88888888888888888880000000000000000000111111111111 ACCEPTING PROMISE: <SELF,TERM,PROSAL_NUM>", a.self, accptMessage.term, accptMessage.proposal_num, accptMessage.value.msgId)
		arg.promised_num = accptMessage.proposal_num
		arg.accepted_num = accptMessage.proposal_num
		arg.accepted_value = accptMessage.value
		arg.remoteAddress = customConn.RemoteAddress().String()
		customConn.SendData2(ACCEPTOR_PROTO_ID, PROPOSER_PROTO_ID, accptMessage, ON_ACCEPTED_ID)
		a.totalSent++
	}
}

// type LocalProtoComHandlerFunc func(sourceProto APP_PROTO_ID, destProto ProtoInterface, data interface{})
func AcceptorValueDecided(sourceProto tcpChannel.APP_PROTO_ID, destProto tcpChannel.ProtoInterface, data interface{}) {
	_, valid := destProto.(*AcceptorProto)
	if valid {
		_, _ = data.(*PaxosMsg)
		/*
			if ok {
				arg := a.terms[decided.term]
				if arg == nil {
					arg = &termArg{}
					a.terms[decided.term] = arg
				}
				arg.decided = decided
			}*/
	}
}
func (a *AcceptorProto) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return ACCEPTOR_PROTO_ID
}
func (a *AcceptorProto) OnStart(channelInterface tcpChannel.ChannelInterface) {
	err1 := a.protoManager.RegisterNetworkMessageHandler(ON_PREPARE_ID, a.onPrepare) //registar no server
	err2 := a.protoManager.RegisterNetworkMessageHandler(ON_ACCEPT_ID, a.onAccept)   //registar no server
	log.Println("REGISTER MSG HANDLERS ACCEPTOR PROTO: ", err1, err2)
	a.protoManager.RegisterPeriodicTimeout(a.ProtocolUniqueId(), time.Second*300, nil, a.PeriodicTimerHandler)
}
func (a *AcceptorProto) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *AcceptorProto) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	a.peers[customCon.RemoteAddress().String()] = customCon

}
func (a *AcceptorProto) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	log.Println("CONNECTION IS DOWN ", a.self)
}
func (c *AcceptorProto) PeriodicTimerHandler(handlerId int, proto tcpChannel.APP_PROTO_ID, message interface{}) {
	log.Println("TIMER ON TRIGGERED ACCEPTOR", c.self, c.totalSent)
}
