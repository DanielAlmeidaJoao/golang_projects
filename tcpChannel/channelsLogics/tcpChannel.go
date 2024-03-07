package channelsLogics

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fileServer/protocolLIstenerLogics"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
)

const MAX_BYTES_TO_READ = 2 * 1024

func abort(err error) bool {
	if err == nil {
		return false
	} else {
		panic(err)
		return true
	}
}

type CONNECTION_TYPE int8

const (
	P2P CONNECTION_TYPE = iota
	CLIENT
	SERVER
)

type MSG_TYPE uint8

var NOT_CONNECTED = errors.New("NOT_CONNECTED")
var UNEXPECTED_NETWORK_READ_DATA = errors.New("UNEXPECTED_NETWORK_READ_DATA")

const (
	LISTEN_ADDRESS_MSG MSG_TYPE = iota
	APP_MSG
)

type TCPChannel struct {
	mu             sync.Mutex
	connections    map[string]net.Conn
	address        string
	port           int
	listener       net.Listener
	protoListener  protocolLIstenerLogics.ProtoListener
	connectionType CONNECTION_TYPE
}

func NewTCPChannel(address string, port int, protoListener protocolLIstenerLogics.ProtoListener, connectionType CONNECTION_TYPE) *TCPChannel {
	channel := &TCPChannel{
		connections:    make(map[string]net.Conn),
		address:        address,
		port:           port,
		protoListener:  protoListener,
		connectionType: connectionType,
	}
	channel.start()
	return channel
}
func shutDown(c *TCPChannel) {
	c.listener.Close()
}
func (c *TCPChannel) formatAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.port)
}
func (c *TCPChannel) start() {
	listener, err := net.Listen("tcp", c.formatAddress())
	abort(err)
	c.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			abort(err)
			c.onConnected(conn, fmt.Sprintf("%s:%d", c.address, c.port), protocolLIstenerLogics.IN_CONNECTION_UP)
		}
	}()
}
func ToByteNetworkData(msgType MSG_TYPE, data []byte) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, msgType)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		os.Exit(1)
	}
	buf.Write(data)
	//netData := make([] byte,1+len(data))
	return buf.Bytes()
}

func (c *TCPChannel) toByteMSG(msgType MSG_TYPE, data []byte) {
	//TODO sum := make([]byte, 1+len(data))
	//append(sum, msgType...)
}
func (c *TCPChannel) onConnected(conn net.Conn, listenAddress string, inConnection protocolLIstenerLogics.ConState) {
	c.mu.Lock()
	c.connections[conn.RemoteAddr().String()] = conn
	c.mu.Unlock()
	//empty string, the user only wants to be a client
	if strings.Compare("", listenAddress) == 0 {
		//TODO SEND LISTEN ADDRESS
		// LATER CHECK IF IT IS CLIENT | SERVER | P2P
	}
	//TODO should all protocols receive connection up event ??
	connectionUp := protocolLIstenerLogics.NewNetworkEvent(conn.RemoteAddr(), nil, -1, inConnection, -1)
	c.protoListener.DeliverEvent(connectionUp)
	//the protocols receive the connection when connection is up
	c.readFromConnection(conn)
}
func (c *TCPChannel) closeConnection(connectionId string) {
	conn := c.connections[connectionId]
	go func() {
		if conn != nil {
			c.mu.Lock()
			delete(c.connections, connectionId)
			c.mu.Unlock()
			conn.Close()
		}
	}()
}
func (c *TCPChannel) closeConnection2(conn net.Conn) {
	c.closeConnection(conn.RemoteAddr().String())
}
func (c *TCPChannel) connect(address string, port int, listenAddress string) {
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			//TODO notify about the error
			return
		}
		c.onConnected(conn, listenAddress, protocolLIstenerLogics.OUT_CONNECTION_UP)
	}()
}

func (c *TCPChannel) readFromConnection(conn net.Conn) {
	for {
		var datalen int
		var err error
		err = binary.Read(conn, binary.LittleEndian, &datalen)
		if connectionDown(err) {

		}
		var msgCode MSG_TYPE
		err = binary.Read(conn, binary.LittleEndian, &msgCode)
		if connectionDown(err) {

		}
		buffer := make([]byte, datalen-1)
		datalen, err = conn.Read(buffer)

		if connectionDown(err) {
			println("CLIENT CONNECTION IS DOWN!")
			return
		}
		//TODO deliver event
	}
}

func (c *TCPChannel) sendMessage(ipAddress string, msg []byte, msgType MSG_TYPE) (int, error) {
	//I DONT LIKE IT:
	// using protoBuf to binary a struct with the data, msgType, ??
	netData := ToByteNetworkData(msgType, msg)
	conn := c.connections[ipAddress]
	written := -1
	var err error
	if conn != nil {
		written, err = conn.Write(netData)
		connectionDown(err)
		return written, err
		//written to logger service
	} else {
		err = NOT_CONNECTED
	}

	return written, err
}

func connectionDown(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)
}
