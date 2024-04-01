package paxos

import (
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
)

const (
	PROPOSER_PROTO_default tcpChannel.APP_PROTO_ID = iota
	PROPOSER_PROTO_ID
	ACCEPTOR_PROTO_ID
	LEARNER_PROTO_ID
	CLIENT_PROTO_ID
)
const (
	ON_PROPOSE_ID tcpChannel.MessageHandlerID = iota
	ON_PROMISE_ID
	ON_ACCEPTED_ID
	ON_ACCEPT_ID
	ON_PREPARE_ID
	ON_DECIDE_ID
)

type PaxosMsg struct {
	msgId    string
	msgValue string
}

func (receiver *PaxosMsg) SerializeData(writer *tcpChannel.CustomWriter) {
	receiver.WritePaxosMsg(writer)
}
func (p *PaxosMsg) WritePaxosMsg(writer *tcpChannel.CustomWriter) {
	//TODO check when the message is NIL
	if p == nil {
		p = &PaxosMsg{}
	}
	writer.WriteUInt32(uint32(len(p.msgId)))
	if len(p.msgId) > 0 {
		writer.WriteString(p.msgId)
		writer.WriteUInt32(uint32(len(p.msgValue)))
		writer.WriteString(p.msgValue)
	}
}

//

func ReadPaxosMsg(data *tcpChannel.CustomReader) *PaxosMsg {
	//TODO check when the message is NIL
	msg := &PaxosMsg{}
	msgIdLen, _ := data.ReadUint32()
	if msgIdLen > 0 {
		bytes := make([]byte, msgIdLen)
		data.Read(bytes)
		msg.msgId = string(bytes)
		msgValueLen, _ := data.ReadUint32()
		bytes = make([]byte, msgValueLen)
		data.Read(bytes)
		msg.msgValue = string(bytes)
	} else {
		msg = nil
	}
	return msg
}

type PrepareMessage struct {
	proposal_num uint32
}

func (receiver *PrepareMessage) SerializeData(writer *tcpChannel.CustomWriter) {
	writer.WriteUInt32(receiver.proposal_num)
}
func ReadDataPrepareMessage(reader *tcpChannel.CustomReader) *PrepareMessage {
	accepted_num, _ := reader.ReadUint32()
	return &PrepareMessage{
		proposal_num: accepted_num,
	}
}

type AcceptMessage struct {
	proposal_num uint32
	value        *PaxosMsg
}

func (receiver *AcceptMessage) SerializeData(writer *tcpChannel.CustomWriter) {
	writer.WriteUInt32(receiver.proposal_num)
	receiver.value.WritePaxosMsg(writer)
}
func ReadDataAcceptMessage(reader *tcpChannel.CustomReader) *AcceptMessage {
	accepted_num, _ := reader.ReadUint32()
	msg := ReadPaxosMsg(reader)
	return &AcceptMessage{
		proposal_num: accepted_num,
		value:        msg,
	}
}

type Promise struct {
	accepted_num   uint32
	promised_num   uint32
	accepted_value *PaxosMsg
}

func (receiver *Promise) SerializeData(writer *tcpChannel.CustomWriter) {
	writer.WriteUInt32(receiver.accepted_num)
	receiver.accepted_value.WritePaxosMsg(writer)
}

func ReadData(reader *tcpChannel.CustomReader) *Promise {
	accepted_num, _ := reader.ReadUint32()
	msg := ReadPaxosMsg(reader)
	return &Promise{
		accepted_num:   accepted_num,
		accepted_value: msg,
	}
}
