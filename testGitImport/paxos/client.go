package paxos

import "C"
import (
	list2 "container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	ops          *list2.List
	count        int
	protocolAPI  tcpChannel.ProtocolAPI
	proper       *ProposerProtocol
	currentTerm  uint32
	lastProposed *PaxosMsg
	self         string
	start        int64
}

func NewClientProtocol(proposer *ProposerProtocol, address string) *ClientProtocol {
	return &ClientProtocol{
		ops:         list2.New(),
		count:       0,
		proper:      proposer,
		currentTerm: uint32(1),
		self:        address,
	}
}
func (c *ClientProtocol) ProtocolUniqueId() tcpChannel.APP_PROTO_ID {
	return CLIENT_PROTO_ID
}
func SendPaxosRequest(sourceProto tcpChannel.APP_PROTO_ID, destProto tcpChannel.ProtoInterface, data interface{}) {
	p, ok := destProto.(*ProposerProtocol)
	if ok {
		paxosMsg, ok := data.(*PaxosMsg)
		if ok {
			p.OnProposeClientCall(paxosMsg)
		}
	} else {
		log.Println("ERROR CONVERTING")
	}
}
func (c *ClientProtocol) nextProposal() *PaxosMsg {
	if c.start == 0 {
		c.start = time.Now().UnixMilli()
	}
	c.count++
	if c.count > 1000 {
		return nil
	}
	return &PaxosMsg{
		msgId:    fmt.Sprintf("%s_%d", c.self, c.count),
		msgValue: fmt.Sprintf("MSG_IS_%s_%d", c.self, c.count),
		term:     c.currentTerm,
	}
}
func (c *ClientProtocol) handleTimer(handlerId int, sourceProto tcpChannel.APP_PROTO_ID, data interface{}) {
	msg := c.nextProposal()
	c.lastProposed = msg
	log.Println("----------------------------------------------------------------------------------------------- timer triggereed")
	f := SendPaxosRequest
	err2 := c.protocolAPI.SendLocalEvent(PROPOSER_PROTO_ID, msg, f) //registar no server
	log.Println("ERROR REGISTERING PROTO", err2)
}

func (c *ClientProtocol) appendMap() string {
	aux := ""
	size := c.ops.Len()
	fmt.Println(c.self, " <", size, ">")

	if size > 0 {
		first := c.ops.Front()
		h := sha256.New()
		for i := 0; i < size; i++ {
			value, ok := first.Value.(*PaxosMsg)
			if ok {
				h.Write([]byte(value.msgValue))
				//aux = fmt.Sprintf("%s <VALUE: %s TERM: %d ; PROPOSAL NUM: %d>", aux, value.msgValue, value.term, value.proposalNum)
			}
			first = first.Next()
		}
		aux = fmt.Sprintf("SERVER: %s -- SIZE: %d -- HASH: %s -- VALUES: %s", c.self, size, hex.EncodeToString(h.Sum(nil)), aux)
	}
	return aux
}
func (c *ClientProtocol) PeriodicTimerHandler(handlerId int, proto tcpChannel.APP_PROTO_ID, message interface{}) {
	log.Println("MAP ----- ++++ ", c.self, c.appendMap())
}

// type LocalProtoComHandlerFunc func(sourceProto APP_PROTO_ID, destProto ProtoInterface, data interface{})
func ValueDecided(sourceProto tcpChannel.APP_PROTO_ID, destProto tcpChannel.ProtoInterface, data interface{}) {
	l := list2.New()
	l.PushBack(12)
	l.Len()
	c, valid := destProto.(*ClientProtocol)
	if valid {

		value, ok := data.(*PaxosMsg)
		if ok {
			if c.currentTerm == value.term {
				c.currentTerm += 1
			}
			c.ops.PushBack(value)
			if c.ops.Len() == 3000 {
				log.Println(c.self, " -- ELAPSED IS -- : ", time.Now().UnixMilli()-c.start)
			}
			if c.lastProposed != nil && value.msgId == c.lastProposed.msgId {
				c.lastProposed = c.nextProposal()
			}
			if c.lastProposed != nil {
				c.lastProposed.term = c.currentTerm
				c.lastProposed.proposalNum = value.proposalNum
				_ = c.protocolAPI.SendLocalEvent(PROPOSER_PROTO_ID, c.lastProposed, SendPaxosRequest) //registar no server
			}

			if c.lastProposed == nil {
				_ = c.protocolAPI.SendLocalEvent(PROPOSER_PROTO_ID, &PaxosMsg{term: c.currentTerm, msgId: "", proposalNum: value.proposalNum}, SendPaxosRequest)
			}

			//log.Println("MAP ----- ++++ ", c.self, c.appendMap())

		}
	}
}

func (c *ClientProtocol) OnStart(protocolAPI tcpChannel.ProtocolAPI) {
	c.protocolAPI = protocolAPI
	time.Sleep(time.Second * 5)
	protocolAPI.NetworkInterface().OpenConnection("127.0.0.1", 8080, c.ProtocolUniqueId())
	protocolAPI.NetworkInterface().OpenConnection("127.0.0.1", 8081, c.ProtocolUniqueId())
	protocolAPI.NetworkInterface().OpenConnection("127.0.0.1", 8082, c.ProtocolUniqueId())
	protocolAPI.RegisterTimeout(time.Second*5, nil, c.handleTimer)
	timerId, err := protocolAPI.RegisterPeriodicTimeout(time.Second*30, nil, c.PeriodicTimerHandler)
	log.Println("REGISTERED TIME ID IS : ", timerId, err)
}
func (c *ClientProtocol) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ProtocolAPI) {

}
func (c *ClientProtocol) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ProtocolAPI) {
	log.Println("CONNECTION IS UPP")
}
func (c *ClientProtocol) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ProtocolAPI) {
	log.Println("CONNECTION IS down")

}
