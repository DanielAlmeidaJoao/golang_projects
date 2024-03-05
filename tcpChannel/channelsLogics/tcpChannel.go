package channelsLogics

import (
	"errors"
	"fileServer/protocolLIstenerLogics"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"syscall"
)

const MAX_BYTES_TO_READ = 64 * 1024

func abort(err error) bool {
	if err == nil {
		return false
	} else {
		panic(err)
		return true
	}
}

type TCPChannel struct {
	mu            sync.Mutex
	connections   map[string]net.Conn
	address       string
	port          int
	listener      net.Listener
	protoListener protocolLIstenerLogics.ProtoListener
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
func (c *TCPChannel) onConnected(conn net.Conn, listenAddress string, inConnection protocolLIstenerLogics.ConState) {
	c.mu.Lock()
	c.connections[conn.RemoteAddr().String()] = conn
	c.mu.Unlock()
	if strings.Compare("", listenAddress) == 0 {
		//TODO SEND LISTEN ADDRESS
		// LATER CHECK IF IT IS CLIENT | SERVER | P2P
	}
	//TODO should all protocols receive connection up event ??
	connectionUp := protocolLIstenerLogics.NewNetworkEvent(conn.RemoteAddr(), nil, -1, inConnection, -1)
	c.protoListener.DeliverEvent(connectionUp)
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
		buffer := make([]byte, MAX_BYTES_TO_READ)
		_, err := conn.Read(buffer)
		if connectionDown(err) {
			println("CLIENT CONNECTION IS DOWN!")
			return
		}
	}
}

func (c *TCPChannel) sendMessage(ipAddress string, msg []byte) bool {
	conn := c.connections[ipAddress]
	if conn != nil {
		written, err := conn.Write(msg)
		connectionDown(err)
		return written > 0
		//written to logger service
	}
	return false
}

func connectionDown(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)
}
