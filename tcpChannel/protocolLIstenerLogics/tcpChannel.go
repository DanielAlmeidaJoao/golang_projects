package protocolLIstenerLogics

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	gobabelUtils "gobabel/commons"
	"io"
	"net"
	"strings"
	"sync"
	"syscall"
)

type ChannelInterface interface {
	SendAppData(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error)
	CloseConnection(ipAddress string)
	OpenConnection(address string, port int, listenAddress string, protoSource gobabelUtils.APP_PROTO_ID)
}

func NewTCPChannel(address string, port int, connectionType gobabelUtils.CONNECTION_TYPE) *TCPChannel {
	channel := &TCPChannel{
		connections:    make(map[string]net.Conn),
		address:        address,
		port:           port,
		connectionType: connectionType,
	}
	channel.start()
	return channel
}

type TCPChannel struct {
	mutex          sync.Mutex
	connections    map[string]net.Conn
	address        string
	port           int
	listener       net.Listener
	protoListener  ProtoListener
	connectionType gobabelUtils.CONNECTION_TYPE
}

func shutDown(c *TCPChannel) {
	c.listener.Close()
}
func (c *TCPChannel) SetProtoLister(p *ProtoListener) {
	c.protoListener = *p
}
func (c *TCPChannel) formatAddress() string {
	return fmt.Sprintf("%s:%d", c.address, c.port)
}
func (c *TCPChannel) start() {
	listener, err := net.Listen("tcp", c.formatAddress())
	gobabelUtils.Abort(err)
	c.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			gobabelUtils.Abort(err)
			c.onConnected(conn, fmt.Sprintf("%s:%d", c.address, c.port), gobabelUtils.ALL_PROTO_ID, true)
		}
	}()
}

func (c *TCPChannel) toByteMSG(msgType gobabelUtils.MSG_TYPE, data []byte) {
	//TODO sum := make([]byte, 1+len(data))
	//append(sum, msgType...)
}
func (c *TCPChannel) onConnected(conn net.Conn, listenAddress string, protoDest gobabelUtils.APP_PROTO_ID, incoming bool) {
	c.mutex.Lock()
	c.connections[conn.RemoteAddr().String()] = conn
	c.mutex.Unlock()
	//empty string, the user only wants to be a client
	if strings.Compare("", listenAddress) == 0 {
		//TODO SEND LISTEN ADDRESS
		// LATER CHECK IF IT IS CLIENT | SERVER | P2P
	}
	connectionUp := gobabelUtils.NewNetworkEvent(conn.RemoteAddr(), nil, protoDest, protoDest, gobabelUtils.CONNECTION_UP, -1)
	c.protoListener.DeliverEvent(connectionUp)
	//the protocols receive the connection when connection is up
	c.readFromConnection(&conn)
}
func (c *TCPChannel) CloseConnection(connectionId string) {
	conn := c.connections[connectionId]
	go func() {
		if conn != nil {
			c.mutex.Lock()
			delete(c.connections, connectionId)
			c.mutex.Unlock()
			conn.Close()
		}
	}()
}
func (c *TCPChannel) closeConnection2(conn net.Conn) {
	c.CloseConnection(conn.RemoteAddr().String())
}
func (c *TCPChannel) OpenConnection(address string, port int, listenAddress string, protoSource gobabelUtils.APP_PROTO_ID) {
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			//TODO notify about the error
			return
		}
		c.onConnected(conn, listenAddress, protoSource, false)
	}()
}

func (t *TCPChannel) handleConnectionDown(conn *net.Conn, err *error) {
	//TODO
	(*conn).Close()
}

//write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData

func (c *TCPChannel) readFromConnection(aux *net.Conn) {
	conn := *aux
	for {

		var err error
		var msgCode gobabelUtils.MSG_TYPE
		err = binary.Read(conn, binary.LittleEndian, &msgCode)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(aux, &err)
			return
		}
		var sourceProto uint16
		err = binary.Read(conn, binary.LittleEndian, &sourceProto)
		if connectionDown(&err) {
			c.handleConnectionDown(aux, &err)
			return
		}

		var destProto uint16
		err = binary.Read(conn, binary.LittleEndian, &destProto)
		if connectionDown(&err) {
			c.handleConnectionDown(aux, &err)
			return
		}

		var datalen uint32
		err = binary.Read(conn, binary.LittleEndian, &datalen)
		if connectionDown(&err) {
			c.handleConnectionDown(aux, &err)
			return
		}
		buffer := make([]byte, datalen)
		var readData int
		readData, err = conn.Read(buffer)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(aux, &err)
			return
		}
		if uint32(readData) != datalen {
			err = errors.New(fmt.Sprintf("EXPECTED TO READ %d; BUT ONLY READ %d. CONN CLOSED", datalen, readData))
			c.handleConnectionDown(aux, &err)
			return
		}
		msgEvent := gobabelUtils.NewNetworkEvent(conn.RemoteAddr(), buffer, gobabelUtils.APP_PROTO_ID(sourceProto), gobabelUtils.APP_PROTO_ID(destProto), gobabelUtils.MESSAGE_RECEIVED, -1)
		c.protoListener.DeliverEvent(msgEvent)
		//TODO deliver event
	}
}

func (c *TCPChannel) SendAppData(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error) {
	return c.sendMessage(ipAddress, source, destProto, msg, gobabelUtils.APP_MSG)
}

//TODO use a golang routine just to send messages
//write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData
func (c *TCPChannel) sendMessage(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, msgType gobabelUtils.MSG_TYPE) (int, error) {
	//I DONT LIKE IT:
	// using protoBuf to binary a struct with the data, msgType, ??
	//TODO implement my bytebuffer
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msgType)
	binary.Write(buf, binary.LittleEndian, source)
	binary.Write(buf, binary.LittleEndian, destProto)
	binary.Write(buf, binary.LittleEndian, len(msg))
	binary.Write(buf, binary.LittleEndian, msg)

	conn := c.connections[ipAddress]
	written := -1
	var err error
	if conn != nil {
		written, err = conn.Write(buf.Bytes())
		if connectionDown(&err) {
			c.handleConnectionDown(&conn, &err)
		}
		return written, err
		//written to logger service
	} else {
		err = gobabelUtils.NOT_CONNECTED
	}

	return written, err
}

func connectionDown(err *error) bool {
	return errors.Is(*err, io.EOF) || errors.Is(*err, syscall.EPIPE)
}
