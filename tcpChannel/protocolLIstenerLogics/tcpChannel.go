package protocolLIstenerLogics

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	gobabelUtils "gobabel/commons"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"syscall"
)

const HeaderSize = 11

type ChannelInterface interface {
	SendAppData(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error)
	CloseConnection(ipAddress string)
	OpenConnection(address string, port int, protoSource gobabelUtils.APP_PROTO_ID)
	IsConnected(address string) bool
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
			c.readFromConnection(&conn, nil)
			//c.onConnected(conn, fmt.Sprintf("%s:%d", c.address, c.port), gobabelUtils.ALL_PROTO_ID, true)
		}
	}()
}

func (c *TCPChannel) toByteMSG(msgType gobabelUtils.MSG_TYPE, data []byte) {
	//TODO sum := make([]byte, 1+len(data))
	//append(sum, msgType...)
}
func (c *TCPChannel) onDisconnected(from string) {
	c.mutex.Lock()
	delete(c.connections, from)
	c.mutex.Unlock()
}
func (c *TCPChannel) onConnected(conn net.Conn, listenAddress string, protoDest gobabelUtils.APP_PROTO_ID) {
	//senf the local listen address to the remote server
	c.mutex.Lock()
	c.connections[conn.RemoteAddr().String()] = conn
	c.mutex.Unlock()
	c.sendMessage(conn.RemoteAddr().String(), gobabelUtils.ALL_PROTO_ID, gobabelUtils.ALL_PROTO_ID, []byte(listenAddress), gobabelUtils.LISTEN_ADDRESS_MSG, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
	connectionUp := gobabelUtils.NewNetworkEvent(conn.RemoteAddr(), nil, protoDest, protoDest, gobabelUtils.CONNECTION_UP, -1)
	c.protoListener.DeliverEvent(connectionUp)
	//the protocols receive the connection when connection is up
	c.readFromConnection(&conn, conn.RemoteAddr())
	//empty string, the user only wants to be a client
	if strings.Compare("", listenAddress) == 0 {
		//TODO SEND LISTEN ADDRESS
		// LATER CHECK IF IT IS CLIENT | SERVER | P2P
	}
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
func (c *TCPChannel) OpenConnection(address string, port int, protoSource gobabelUtils.APP_PROTO_ID) {
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
		if err != nil {
			//TODO notify about the error
			return
		}
		c.onConnected(conn, fmt.Sprintf("%s:%d", c.address, c.port), protoSource)
	}()
}
func (c *TCPChannel) readFromConnection2(aux *net.Conn, listenAddress net.Addr) {

}
func writeHeaders(buf io.Writer, source, destProto gobabelUtils.APP_PROTO_ID, msgType gobabelUtils.MSG_TYPE, msgHandlerId gobabelUtils.MessageHandlerID) {
	binary.Write(buf, binary.LittleEndian, uint8(msgType))
	binary.Write(buf, binary.LittleEndian, uint16(msgHandlerId))
	binary.Write(buf, binary.LittleEndian, uint16(source))
	binary.Write(buf, binary.LittleEndian, uint16(destProto))
	buf.Write([]byte{0, 0, 0, 0})
	//binary.Write(buf, binary.LittleEndian, uint32(0))
}

//write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData

func (c *TCPChannel) readFromConnection(aux *net.Conn, listenAddress net.Addr) {
	conn := *aux
	for {
		var err error
		var msgCode uint8
		err = binary.Read(conn, binary.LittleEndian, &msgCode)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(aux, &err)
			return
		}
		var msgHandlerId uint16
		err = binary.Read(conn, binary.LittleEndian, &msgHandlerId)
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
		var dataLen uint32
		err = binary.Read(conn, binary.LittleEndian, &dataLen)
		if connectionDown(&err) {
			c.handleConnectionDown(aux, &err)
			return
		}
		buffer := make([]byte, dataLen)
		var readData int
		readData, err = conn.Read(buffer)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(aux, &err)
			return
		}
		if readData != int(dataLen) {
			str := fmt.Sprintf("EXPECTED TO READ %d; BUT ONLY READ %d. CONN CLOSED", dataLen, readData)
			err = errors.New(str)
			fmt.Println(str)
			c.handleConnectionDown(aux, &err)
			return
		}
		netEvent := gobabelUtils.MESSAGE_RECEIVED
		if listenAddress == nil {
			tcpAddr, err := net.ResolveTCPAddr("tcp", string(buffer))
			if err == nil {
				log.Println("THE CLIENT REMOTE ADDRESS IS: ", tcpAddr.String())
				c.mutex.Lock()
				c.connections[tcpAddr.String()] = conn
				c.mutex.Unlock()
				listenAddress = tcpAddr
				netEvent = gobabelUtils.CONNECTION_UP
			} else {
				log.Fatal("FAILED TO PARSE RECEIVED ADDRESS OF THE SERVER:", err)
				conn.Close()
				return
			}
			buffer = nil
		}
		msgEvent := gobabelUtils.NewNetworkEvent(listenAddress, buffer, gobabelUtils.APP_PROTO_ID(sourceProto), gobabelUtils.APP_PROTO_ID(destProto), netEvent, gobabelUtils.MessageHandlerID(msgHandlerId))
		c.protoListener.DeliverEvent(msgEvent)
	}
}

func (c *TCPChannel) SendAppData(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error) {
	return c.sendMessage(ipAddress, source, destProto, msg, gobabelUtils.APP_MSG, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
}
func (c *TCPChannel) SendAppData2(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	buf := new(bytes.Buffer)
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	writeHeaders(buf, source, destProto, gobabelUtils.APP_MSG, msgHandlerId)
	msg.serializeData(NewBinaryWriter2(buf, binary.LittleEndian))
	return c.auxWriteToNetwork(hostAddress, buf.Bytes())
}

func (c *TCPChannel) IsConnected(address string) bool {
	return c.connections[address] != nil
}

func (c *TCPChannel) auxWriteToNetwork(address string, data []byte) (int, error) {
	conn := c.connections[address]
	written := -1
	var err error
	if conn != nil {
		dataSize := uint32(len(data) - HeaderSize)
		writer := NewCustomWriter2(data, HeaderSize)
		_ = binary.Write(writer, binary.LittleEndian, dataSize)
		written, err = conn.Write(data)
		if connectionDown(&err) {
			c.handleConnectionDown(&conn, &err)
		}
		return written, err
		//written to logger service
	} else {
		err = gobabelUtils.NOT_CONNECTED
	}
	fmt.Println("GOING TO SEND MESSAGE")
	return written, err
}

// TODO use a golang routine just to send messages
// write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData
func (c *TCPChannel) sendMessage(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, msgType gobabelUtils.MSG_TYPE, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	//I DONT LIKE IT:
	// using protoBuf to binary a struct with the data, msgType, ??
	fmt.Println("FIRST MESSAGE IS: ", string(msg))
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	buf := new(bytes.Buffer)
	writeHeaders(buf, source, destProto, msgType, msgHandlerId)
	err := binary.Write(buf, binary.LittleEndian, msg)
	if err != nil {
		log.Fatal("FAILED TO APPEND THE BINARY DATA!", err)
		return 0, err
	}
	return c.auxWriteToNetwork(hostAddress, buf.Bytes())
}
func (t *TCPChannel) handleConnectionDown(conn *net.Conn, err *error) {
	fmt.Println("RECEIVED MESSAGE. DELIVERING EVENTS ERRROR ", *err)
	(*conn).Close()
	t.onDisconnected((*conn).RemoteAddr().String())
	msgEvent := gobabelUtils.NewNetworkEvent((*conn).RemoteAddr(), nil, gobabelUtils.ALL_PROTO_ID, gobabelUtils.ALL_PROTO_ID, gobabelUtils.CONNECTION_DOWN, -1)
	t.protoListener.DeliverEvent(msgEvent)
}
func connectionDown(err *error) bool {
	return errors.Is(*err, io.EOF) || errors.Is(*err, syscall.EPIPE)
}
