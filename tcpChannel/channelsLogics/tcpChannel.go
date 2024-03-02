package channelsLogics

import (
	"errors"
	"fmt"
	"io"
	"net"
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
	mu          sync.Mutex
	connections map[string]net.Conn
	address     string
	port        int
	listener    net.Listener
}

func shutDown(c *TCPChannel) {
	c.listener.Close()
}
func (c *TCPChannel) start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))
	abort(err)
	c.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			abort(err)
			c.onConnected(conn)
		}
	}()
}
func (c *TCPChannel) onConnected(conn net.Conn) {
	c.mu.Lock()
	c.connections[conn.RemoteAddr().String()] = conn
	c.mu.Unlock()
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
func (c *TCPChannel) connect(address string, port int) {
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			//TODO notify about the error
			return
		}
		c.onConnected(conn)
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
func connectionDown(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)
}
