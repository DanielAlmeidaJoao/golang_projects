package main

import (
	"encoding/binary"
	"fmt"
	"github.com/DanielAlmeidaJoao/goDistributedLibrary/tcpChannel"
	"log"
	"os"
	"strconv"
	"strings"
	"testGitImport/paxos"
)

// GONOPROXY=github.com/DanielAlmeidaJoao go get -d github.com/DanielAlmeidaJoao/goDistributedLibrary
func main() {
	selfAddress := os.Args[1]
	before, after, _ := strings.Cut(selfAddress, ":")
	log.Println("AFTER is ", after)
	port, _ := strconv.Atoi(after)
	log.Println("LISTENING ON: ", before, port)
	protocolsManager := tcpChannel.NewProtocolListener(before, port, tcpChannel.SERVER, binary.LittleEndian)
	//fmt.Println(protocolsManager)
	fmt.Println(protocolsManager)

	proposer := paxos.NewProposerProtocol(protocolsManager)
	client := paxos.NewClientProtocol(protocolsManager, proposer)
	acceptor := paxos.NewAcceptorProtocol(protocolsManager)
	learner := paxos.NewLearnerProtocol(protocolsManager, client)

	protocolsManager.StartProtocol(proposer)
	protocolsManager.StartProtocol(acceptor)
	protocolsManager.StartProtocol(learner)
	protocolsManager.StartProtocol(client)

	protocolsManager.WaitForProtocolsToEnd(false)

}
