package testUtils

import protoListener "gobabel/protocolLIstenerLogics"

/*** ***/
type EchoMessage struct {
	Data  string
	Count int32
}

func (receiver *EchoMessage) SerializeData(writer *protoListener.CustomWriter) {
	//TODO implement me
	_, _ = writer.WriteUInt32(uint32(receiver.Count))
	_, _ = writer.WriteString(receiver.Data)
}

func DeserializeData(reader *protoListener.CustomReader) *EchoMessage {
	count, _ := reader.ReadUint32()
	str := string(reader.ReadAll())
	aux := &EchoMessage{
		Data:  str,
		Count: int32(count),
	}
	return aux
}
