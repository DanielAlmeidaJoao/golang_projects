package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"math"
)

// the state of the proposer
type ProposerProtocol struct {
	value        *PaxosMsg //proposed value
	proposal_num uint32    //proposal number
	acks         uint32    //number of acks received from acceptors
	promises     int       //promises received from acceptors
	peers        map[string]*tcpChannel.CustomConnection
	ProtoManager tcpChannel.ProtoListenerInterface
	currentTerm  uint32
	self         string
}

func NewProposerProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *ProposerProtocol {
	return &ProposerProtocol{
		ProtoManager: listenerInterface,
		currentTerm:  1,
		self:         address,
		promises:     0,
	}
}

//func (p *ProtoEcho) ClientHandleMessage(customConn *protoListener.CustomConnection, protoSource protoListener.APP_PROTO_ID, data *protoListener.CustomReader) {
/*
   The client calls this function of the proposer to propose its desired value.
   This will be called also after abort to try again.
*/
func (p *ProposerProtocol) OnProposeClientCall(clientValue *PaxosMsg) {
	p.proposal_num++
	//log.Println(p.self, "************************************************** going to propopse with ", clientValue.msgValue, clientValue.term, p.proposal_num)
	p.value = clientValue
	p.acks = 0
	prepMessage := &PrepareMessage{
		proposal_num: p.proposal_num,
		term:         clientValue.term,
	}
	p.currentTerm = clientValue.term
	p.promises = 0
	p.acks = 0
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
	if promise.term < p.currentTerm {
		return
	}
	if promise.accepted_value == nil {
		promise.promised_num = p.proposal_num
		promise.accepted_value = p.value
	}

	if promise.promised_num == p.proposal_num {
		p.promises++
		if p.promises == p.majority() {
			p.acks = 0
			p.promises = 0
			if promise.accepted_num > p.proposal_num {
				p.value = promise.accepted_value
			}
			amsg := &AcceptMessage{
				value:        p.value,
				proposal_num: p.proposal_num,
				term:         promise.term,
			}
			for _, v := range p.peers {
				v.SendData2(PROPOSER_PROTO_ID, ACCEPTOR_PROTO_ID, amsg, ON_ACCEPT_ID)
			}
		}
	}

}

/*
This function is called when the proposer receives an accepted message from one of the acceptors.
*/
func (p *ProposerProtocol) onAccepted(customConn *tcpChannel.CustomConnection, protoSource tcpChannel.APP_PROTO_ID, data *tcpChannel.CustomReader) {
	log.Println("ACCEPTED ACCEPTED")
	accept := ReadDataAcceptMessage(data)
	if accept.term < p.currentTerm {
		return
	}
	if p.proposal_num == accept.proposal_num {
		accept.value.proposalNum = p.proposal_num
		p.acks++
		if int(p.acks) == p.majority() {
			p.currentTerm++
			p.acks = 0
			log.Println(p.self, " ?????? ", accept.value.msgValue, "TERM ", accept.value.term, " PROMISE COUNT", p.promises)
			for _, v := range p.peers {
				v.SendData2(PROPOSER_PROTO_ID, LEARNER_PROTO_ID, accept.value, ON_DECIDE_ID)
			}
		}
	}
}

func (a *ProposerProtocol) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return PROPOSER_PROTO_ID
}
func (a *ProposerProtocol) OnStart(channelInterface tcpChannel.ChannelInterface) {
	//err1 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROPOSE_ID, a.onPropose)   //registar no server
	err2 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROMISE_ID, a.onPromise)   //registar no server
	err3 := a.ProtoManager.RegisterNetworkMessageHandler(ON_ACCEPTED_ID, a.onAccepted) //registar no server

	log.Println("REGISTER MSG HANDLERS: ", err2, err3)
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
