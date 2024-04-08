package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
)

/************************************ ***************************************/
type AcceptorProto struct {
	accepted_num   uint32
	promised_num   uint32
	term           uint32
	accepted_value *PaxosMsg
	protoManager   tcpChannel.ProtoListenerInterface
}

func NewAcceptorProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *AcceptorProto {
	return &AcceptorProto{
		protoManager: listenerInterface,
		term:         1,
	}
}
func (a *AcceptorProto) onPrepare(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	preparaMessage := ReadDataPrepareMessage(data)
	log.Println("ON PREPARE 2 ", a.promised_num, preparaMessage.proposal_num, customConn.RemoteAddress().String())
	if preparaMessage.term < a.term {
		return
	}
	if a.promised_num < preparaMessage.proposal_num {
		a.promised_num = preparaMessage.proposal_num
		prom := &Promise{
			accepted_value: a.accepted_value,
			accepted_num:   a.accepted_num,
			promised_num:   preparaMessage.proposal_num,
			term:           preparaMessage.term,
		}
		customConn.SendData2(ACCEPTOR_PROTO_ID, PROPOSER_PROTO_ID, prom, ON_PROMISE_ID)
	}
}

func (a *AcceptorProto) onAccept(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	log.Println("ACCEPTOR RECEIVED ACCEPT MSG")
	accptMessage := ReadDataAcceptMessage(data)
	if accptMessage.term < a.term {
		return
	}
	if a.promised_num <= accptMessage.proposal_num {
		log.Println("INTO AACCEPTER ON ACCEPT")
		a.promised_num = accptMessage.proposal_num
		a.accepted_num = accptMessage.proposal_num
		a.accepted_value = accptMessage.value
		customConn.SendData2(ACCEPTOR_PROTO_ID, PROPOSER_PROTO_ID, accptMessage, ON_ACCEPTED_ID)
		log.Println("EXITED AACCEPTED ON ACCEPT")

	}
}

// type LocalProtoComHandlerFunc func(sourceProto APP_PROTO_ID, destProto ProtoInterface, data interface{})
func AcceptorValueDecided(sourceProto tcpChannel.APP_PROTO_ID, destProto tcpChannel.ProtoInterface, data interface{}) {
	a, valid := destProto.(*AcceptorProto)
	if valid {
		a.term += 1
	}
	a.accepted_value = nil
}
func (a *AcceptorProto) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return ACCEPTOR_PROTO_ID
}
func (a *AcceptorProto) OnStart(channelInterface tcpChannel.ChannelInterface) {
	err1 := a.protoManager.RegisterNetworkMessageHandler(ON_PREPARE_ID, a.onPrepare) //registar no server
	err2 := a.protoManager.RegisterNetworkMessageHandler(ON_ACCEPT_ID, a.onAccept)   //registar no server
	log.Println("REGISTER MSG HANDLERS ACCEPTOR PROTO: ", err1, err2)
}
func (a *AcceptorProto) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *AcceptorProto) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
func (a *AcceptorProto) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {

}
