package commons

import (
	"bytes"
	"encoding/binary"
	"fileServer/channelsLogics"
	"log"
	"net"
)

type DataWrapper struct {
	From    net.Addr
	Data    []byte
	protoId int16
}

func ConvertDataToBinary(num any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		log.Panicf("FAILED TO CONVERT DATA TO BINARY")
	}
	//netData := make([] byte,1+len(data))
	return buf.Bytes(), err
}

func ReadInt(from []byte, conn net.Conn) {
	var dest int
	err = binary.Read(conn, binary.LittleEndian, &dest)
}

func ToByteNetworkData(msgType channelsLogics.MSG_TYPE, data []byte) []byte {
	msgLen, _ := ConvertDataToBinary(1 + len(data))
	msgTypeBin, _ := ConvertDataToBinary(msgType)
	totalLen := len(msgLen) + len(msgTypeBin) + len(data)
	buf := new(bytes.Buffer)
	buf.Grow(totalLen)
	buf.Write(msgLen)
	buf.Write(msgTypeBin)
	buf.Write(data)
	return buf.Bytes()
}
