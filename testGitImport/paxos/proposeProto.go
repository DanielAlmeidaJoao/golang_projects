package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"math"
	"math/rand"
	"time"
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
	rd           *rand.Rand
	timerId      int
}

func NewProposerProtocol(listenerInterface tcpChannel.ProtoListenerInterface, address string) *ProposerProtocol {
	return &ProposerProtocol{
		ProtoManager: listenerInterface,
		currentTerm:  1,
		self:         address,
		promises:     0,
	}
}

func (p *ProposerProtocol) OnProposeClientCall(clientValue *PaxosMsg) {
	p.currentTerm = clientValue.term
	if clientValue.msgId == "" {
		p.value = nil
		p.ProtoManager.CancelTimer(p.timerId)
		p.timerId = -1
		return
	}
	if p.timerId <= 0 {
		p.timerId = p.ProtoManager.RegisterPeriodicTimeout(p.ProtocolUniqueId(), time.Millisecond*time.Duration(150+p.rd.Intn(500)), nil, p.PeriodicTimerHandler)
	}
	p.proposal_num++ //uint32(1 + p.rd.Intn(time.Now().Nanosecond()))
	clientValue.proposalNum = p.proposal_num
	p.value = clientValue
	prepMessage := &PrepareMessage{
		proposal_num: p.proposal_num,
		term:         clientValue.term,
	}
	p.promises = 0
	p.acks = 0
	//log.Println("GOING TO PROPOSE ............................. <SELF,TERM,PROP_NUM>", p.self, p.currentTerm, p.proposal_num)
	for _, v := range p.peers {
		v.SendData2(PROPOSER_PROTO_ID, ACCEPTOR_PROTO_ID, prepMessage, ON_PREPARE_ID)
	}
}
func (c *ProposerProtocol) PeriodicTimerHandler(handlerId int, proto tcpChannel.APP_PROTO_ID, message interface{}) {
	if c.value != nil {
		//log.Println("TIMER TRIGGERTED 7777777777777777777777777777777777777777777777777777777777777777777777777777")
		c.OnProposeClientCall(c.value)
	} else {
		c.timerId = -1
		c.ProtoManager.CancelTimer(handlerId)
	}
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
	if promise.term != p.currentTerm {
		return
	}
	if promise.accepted_value == nil {
		promise.accepted_value = p.value
	}
	if promise.promised_num == p.proposal_num {
		p.promises++
		//log.Println("9999999999999999999999999999999999999999999999999999 RECEIVED PROMISE REPLY : <SELF,TERM,PROPOSAL_NUM,COUNT>", p.self, p.currentTerm, p.proposal_num, p.promises)
		if p.promises == p.majority() {
			promise.accepted_value.proposalNum = p.proposal_num
			promise.accepted_value.term = p.currentTerm
			p.promises = 0
			amsg := &AcceptMessage{
				value:        promise.accepted_value,
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
	accept := ReadDataAcceptMessage(data)
	if accept.term != p.currentTerm {
		return
	}
	if p.proposal_num == accept.proposal_num {
		accept.value.proposalNum = p.proposal_num
		p.acks++
		if int(p.acks) == p.majority() {
			p.acks = 0
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
	a.rd = rand.New(rand.NewSource(time.Now().UnixNano()))
	//err1 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROPOSE_ID, a.onPropose)   //registar no server
	err2 := a.ProtoManager.RegisterNetworkMessageHandler(ON_PROMISE_ID, a.onPromise) //registar no server
	//err3 := a.ProtoManager.RegisterNetworkMessageHandler(ON_ACCEPTED_ID, a.onAccepted) //registar no server

	log.Println("REGISTER MSG HANDLERS: ", err2)
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
