package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"math"
)

// the state of the proposer
type ProposerProtocol struct {
	value        *PaxosMsg  //proposed value
	proposal_num uint32     //proposal number
	acks         uint32     //number of acks received from acceptors
	promises     []*Promise //promises received from acceptors
	peers        map[string]*tcpChannel.CustomConnection
	ProtoManager tcpChannel.ProtoListenerInterface
}

func NewProposerProtocol(listenerInterface tcpChannel.ProtoListenerInterface) *ProposerProtocol {
	return &ProposerProtocol{
		ProtoManager: listenerInterface,
	}
}

//func (p *ProtoEcho) ClientHandleMessage(customConn *protoListener.CustomConnection, protoSource protoListener.APP_PROTO_ID, data *protoListener.CustomReader) {
/*
   The client calls this function of the proposer to propose its desired value.
   This will be called also after abort to try again.
*/
func (p *ProposerProtocol) OnProposeClientCall(clientValue *PaxosMsg) {
	p.proposal_num = p.proposal_num + 1
	p.promises = make([]*Promise, 3)
	p.value = clientValue
	p.acks = 0
	p.promises = nil
	prepMessage := &PrepareMessage{
		proposal_num: p.proposal_num,
	}
	for _, v := range p.peers {
		v.SendData2(PROPOSER_PROTO_ID, ACCEPTOR_PROTO_ID, prepMessage, ON_PREPARE_ID)
	}
}

func (p *ProposerProtocol) onPropose(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	msg := ReadPaxosMsg(data)
	p.OnProposeClientCall(msg)
}
func (p *ProposerProtocol) MaxValue(promises []*Promise) *Promise {
	aux := promises[0]
	for _, v := range promises {
		if v.accepted_value != nil {
			aux = v
		}
	}
	if aux.accepted_value == nil {
		aux.accepted_value = p.value
	}
	return aux
}
func (p *ProposerProtocol) majority() int {
	return int(math.Ceil(float64((len(p.peers) + 1) / 2)))
}

/*
*
This function is called when the proposer receives a promise message from one of the acceptors.
*/
func (p *ProposerProtocol) onPromise(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	promise := ReadData(data)
	log.Println("ACCEPTED PROMISE from ", customConn.RemoteAddress().String(), p.proposal_num, promise.promised_num)
	if promise.promised_num == 0 {
		promise.promised_num = p.proposal_num
		promise.accepted_value = p.value
	}
	if promise.promised_num == p.proposal_num {
		p.promises = append(p.promises, promise)
		majority := p.majority()
		if len(p.promises) == majority {
			aux := p.MaxValue(p.promises)
			if p.value.msgId != aux.accepted_value.msgId {
				p.acks = 0
			}
			p.value = aux.accepted_value
			p.promises = make([]*Promise, 3)
			amsg := &AcceptMessage{
				value:        p.value,
				proposal_num: p.proposal_num,
			}
			for _, v := range p.peers {
				v.SendData2(PROPOSER_PROTO_ID, ACCEPTOR_PROTO_ID, amsg, ON_ACCEPT_ID)
			}
			p.promises = make([]*Promise, len(p.peers))
		}
	}

}

/*
This function is called when the proposer receives an accepted message from one of the acceptors.
*/
func (p *ProposerProtocol) onAccepted(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	log.Println("ACCEPTED ACCEPTED")
	accept := ReadDataAcceptMessage(data)
	p.promises = make([]*Promise, 3)
	if p.proposal_num == accept.proposal_num {
		p.acks++
		if int(p.acks) == p.majority() {
			for _, v := range p.peers {
				v.SendData2(PROPOSER_PROTO_ID, LEARNER_PROTO_ID, accept.value, ON_DECIDE_ID)
			}
		}
	}
}

// the accepter should just ignore it
func (p *ProposerProtocol) onNack(proposal_num uint32) {
	if p.proposal_num == proposal_num {
		//abort
		p.proposal_num = 0
	}
}

func (a *ProposerProtocol) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return PROPOSER_PROTO_ID
}
func (a *ProposerProtocol) OnStart(channelInterface tcpChannel.ChannelInterface) {
	err1 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROPOSE_ID, a.onPropose)   //registar no server
	err2 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROMISE_ID, a.onPromise)   //registar no server
	err3 := a.ProtoManager.RegisterNetworkMessageHandler(ON_ACCEPTED_ID, a.onAccepted) //registar no server

	log.Println("REGISTER MSG HANDLERS: ", err1, err2, err3)
}
func (a *ProposerProtocol) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (a *ProposerProtocol) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	if a.peers == nil {
		a.peers = make(map[string]*tcpChannel.CustomConnection)
	}
	a.peers[customCon.RemoteAddress().String()] = customCon
	log.Println("CONNECTION IS UP", customCon.GetConnectionKey(), customCon.RemoteAddress())
}
func (a *ProposerProtocol) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	log.Println("CONNECTION DOWN", customCon.GetConnectionKey(), customCon.RemoteAddress())
}
