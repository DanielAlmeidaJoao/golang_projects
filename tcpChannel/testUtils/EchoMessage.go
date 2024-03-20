package testUtils

import protoListener "gobabel/protocolLIstenerLogics"

/*** ***/
type EchoMessage struct {
	Data  string
	Count int32
}

func (receiver *EchoMessage) SerializeData(writer *protoListener.CustomWriter) {
	//TODO implement me
	writer.WriteNumber(uint32(receiver.Count))
	writer.WriteString(receiver.Data)
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
