package main

import (
	"encoding/binary"
	"fmt"
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
)

//GONOPROXY=github.com/DanielAlmeidaJoao go get -d github.com/DanielAlmeidaJoao/goDistributedLibrary
func main() {
	protocolsManager := tcpChannel.NewProtocolListener("localhost", 3002, tcpChannel.SERVER, binary.LittleEndian)
	//fmt.Println(protocolsManager)
	fmt.Println(protocolsManager)

}
