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

	proposer := paxos.NewProposerProtocol(selfAddress)
	client := paxos.NewClientProtocol(proposer, selfAddress)
	acceptor := paxos.NewAcceptorProtocol(selfAddress)
	learner := paxos.NewLearnerProtocol(selfAddress)

	//networkQueueSize int, timeoutQueueSize int, localCommQueueSize int
	protocolsManager.StartProtocol(proposer, 100, 50, 20)
	protocolsManager.StartProtocol(acceptor, 100, 50, 20)
	protocolsManager.StartProtocol(learner, 100, 50, 20)
	protocolsManager.StartProtocol(client, 20, 20, 20)

	protocolsManager.WaitForProtocolsToEnd(false)

}
