package paxos

import (
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
	ops          map[string]*PaxosMsg
	count        int
	protoManager tcpChannel.ProtoListenerInterface
	proper       *ProposerProtocol
	global_count uint32
	currentTerm  uint32
	lastProposed *PaxosMsg
	self         string
}

func NewClientProtocol(listenerInterface tcpChannel.ProtoListenerInterface, proposer *ProposerProtocol, address string) *ClientProtocol {
	return &ClientProtocol{
		ops:          make(map[string]*PaxosMsg),
		protoManager: listenerInterface,
		count:        0,
		proper:       proposer,
		currentTerm:  uint32(1),
		self:         address,
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
func (c *ClientProtocol) handleTimer(sourceProto tcpChannel.APP_PROTO_ID, data interface{}) {
	if c.ops == nil {
		c.ops = make(map[string]*PaxosMsg)
	}
	msg := &PaxosMsg{
		msgId:    c.self,
		msgValue: fmt.Sprintf("SELF IS %s", c.self),
		term:     c.currentTerm,
	}
	c.lastProposed = msg
	c.ops[msg.msgId] = msg
	msg.proposalNum = c.global_count + 1
	log.Println("----------------------------------------------------------------------------------------------- timer triggereed")
	f := SendPaxosRequest
	err2 := c.protoManager.SendLocalEvent(c.ProtocolUniqueId(), PROPOSER_PROTO_ID, msg, f) //registar no server
	log.Println("ERROR REGISTERING PROTO", err2)
}
func (c *ClientProtocol) appendMap() string {
	aux := ""
	for _, v := range c.ops {
		aux = fmt.Sprintf("%s ; VALUE: %s TERM: %d ; PROPOSAL NUM: %d ]", aux, v.msgValue, v.term, v.proposalNum)
	}
	return aux
}

// type LocalProtoComHandlerFunc func(sourceProto APP_PROTO_ID, destProto ProtoInterface, data interface{})
func ValueDecided(sourceProto tcpChannel.APP_PROTO_ID, destProto tcpChannel.ProtoInterface, data interface{}) {
	c, valid := destProto.(*ClientProtocol)
	if valid {

		value, ok := data.(*PaxosMsg)
		if ok {
			if c.currentTerm == value.term {
				c.currentTerm += 1
			}
			c.global_count = value.proposalNum
			aux := c.ops[value.msgId]
			/**
			if value.msgId != c.lastProposed.msgId {
				c.lastProposed.term = c.currentTerm
				c.lastProposed.proposalNum = c.global_count
				f := SendPaxosRequest
				_ = c.protoManager.SendLocalEvent(c.ProtocolUniqueId(), PROPOSER_PROTO_ID, c.lastProposed, f) //registar no server

			}**/
			if aux == nil {
				c.ops[value.msgId] = value
				/* TODO FIX THIS N THE GLOBAL COUNTER
				aux = c.GetNonZero()
				f := SendPaxosRequest
				err2 := c.protoManager.SendLocalEvent(c.ProtocolUniqueId(), PROPOSER_PROTO_ID, aux, f) //registar no server
				log.Println("ERROR REGISTERING PROTO", err2)
				*/
			}
			//log.Println("OPS: ", len(c.ops))
			log.Println("MAP ----- ++++ ", c.self, c.appendMap())
		}
	}
}
func (c *ClientProtocol) GetNonZero() *PaxosMsg {
	for _, v := range c.ops {
		if v != nil {
			return v
		}
	}
	return nil
}
func (c *ClientProtocol) OnStart(channelInterface tcpChannel.ChannelInterface) {
	time.Sleep(time.Second * 5)
	channelInterface.OpenConnection("127.0.0.1", 8080, c.ProtocolUniqueId())
	channelInterface.OpenConnection("127.0.0.1", 8081, c.ProtocolUniqueId())
	channelInterface.OpenConnection("127.0.0.1", 8082, c.ProtocolUniqueId())
	c.protoManager.RegisterTimeout(c.ProtocolUniqueId(), time.Second*5, nil, c.handleTimer)
}
func (c *ClientProtocol) OnMessageArrival(customCon *tcpChannel.CustomConnection, source, destProto tcpChannel.APP_PROTO_ID, msg []byte, channelInterface tcpChannel.ChannelInterface) {

}
func (c *ClientProtocol) ConnectionUp(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	log.Println("CONNECTION IS UPP")
}
func (c *ClientProtocol) ConnectionDown(customCon *tcpChannel.CustomConnection, channelInterface tcpChannel.ChannelInterface) {
	log.Println("CONNECTION IS down")

}
