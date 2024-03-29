package protocolLIstenerLogics

import (
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
	SendAppData(addressKey string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error)
	SendAppData2(addressKey string, source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error)

	CloseConnection(ipAddress string)
	CloseConnections()
	OpenConnection(address string, port int, protoSource gobabelUtils.APP_PROTO_ID)
	IsConnected(address string) bool
}

func NewTCPChannel(address string, port int, connectionType gobabelUtils.CONNECTION_TYPE) *TCPChannel {
	channel := &TCPChannel{
		connections:    make(map[string]*CustomConnection),
		address:        address,
		port:           port,
		connectionType: connectionType,
	}
	channel.start()
	return channel
}

type CustomConnection struct {
	conn             net.Conn
	remoteListenAddr net.Addr
	connectionKey    string
}

func (c *CustomConnection) SendData2(source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	customWriter := NewCustomWriter3(binary.LittleEndian)
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	writeHeaders(customWriter, source, destProto, gobabelUtils.APP_MSG, msgHandlerId)
	msg.SerializeData(customWriter)
	return auxWriteToNetworkChild(c, customWriter)
}

func (c *TCPChannel) CloseConnections() {
	for _, v := range c.connections {
		v.conn.Close()
	}
}
func (c *CustomConnection) SendData(source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	customWriter := NewCustomWriter(HeaderSize+len(msg), binary.LittleEndian)
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	writeHeaders(customWriter, source, destProto, gobabelUtils.APP_MSG, msgHandlerId)
	_, _ = customWriter.Write(msg)
	return auxWriteToNetworkChild(c, customWriter)
}

func (c *CustomConnection) RemoteAddress() net.Addr {
	return c.remoteListenAddr
}

func (c *CustomConnection) GetConnectionKey() string {
	return c.connectionKey
}

func (c *CustomConnection) Close() error {
	return c.conn.Close()
}

type TCPChannel struct {
	mutex          sync.Mutex
	connections    map[string]*CustomConnection //map[string]net.Conn  // connectionId, int_ip
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
			customCon := &CustomConnection{
				conn: conn,
			}
			c.readFromConnection(customCon, nil)
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
func (c *TCPChannel) auxAddCon(customCon *CustomConnection, conn net.Conn, listenAddress net.Addr) *CustomConnection {
	if customCon == nil {
		customCon = &CustomConnection{
			conn: conn,
		}
	}

	var addressKey string
	if c.connections[listenAddress.String()] == nil {
		addressKey = listenAddress.String()
	} else {
		addressKey = conn.RemoteAddr().String()
	}
	c.mutex.Lock()
	c.connections[addressKey] = customCon
	c.mutex.Unlock()
	customCon.connectionKey = addressKey
	customCon.remoteListenAddr = listenAddress
	return customCon
}
func (c *TCPChannel) onConnected(conn net.Conn, listenAddress string, protoDest gobabelUtils.APP_PROTO_ID) {
	//send the local listen address to the remote server
	var customCon *CustomConnection
	customCon = c.auxAddCon(nil, conn, conn.RemoteAddr())
	_, err := c.sendMessage(customCon.connectionKey, gobabelUtils.ALL_PROTO_ID, gobabelUtils.ALL_PROTO_ID, []byte(listenAddress), gobabelUtils.LISTEN_ADDRESS_MSG, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
	if err != nil {
		c.handleConnectionDown(customCon, &err)
		return
	}

	connectionUp := NewNetworkEvent(customCon, nil, protoDest, protoDest, gobabelUtils.CONNECTION_UP, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
	c.protoListener.DeliverEvent(connectionUp)
	//the protocols receive the connection when connection is up
	c.readFromConnection(customCon, conn.RemoteAddr())
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
			customCon := &CustomConnection{
				conn: conn, connectionKey: address, remoteListenAddr: nil,
			}
			c.handleConnectionDown(customCon, &err)
		} else {
			c.onConnected(conn, fmt.Sprintf("%s:%d", c.address, c.port), protoSource)
		}
	}()
}
func writeHeaders(buf *CustomWriter, source, destProto gobabelUtils.APP_PROTO_ID, msgType gobabelUtils.MSG_TYPE, msgHandlerId gobabelUtils.MessageHandlerID) {
	_ = binary.Write(buf, binary.LittleEndian, uint8(msgType))
	_ = binary.Write(buf, binary.LittleEndian, uint16(msgHandlerId))
	_ = binary.Write(buf, binary.LittleEndian, uint16(source))
	_ = binary.Write(buf, binary.LittleEndian, uint16(destProto))
	buf.SetOffSet(buf.OffSet() + 4)
}

//write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData

func (c *TCPChannel) readFromConnection(customCon *CustomConnection, listenAddress net.Addr) {
	conn := customCon.conn
	for {
		var err error
		var msgCode uint8
		err = binary.Read(conn, binary.LittleEndian, &msgCode)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(customCon, &err)
			return
		}
		var msgHandlerId uint16
		err = binary.Read(conn, binary.LittleEndian, &msgHandlerId)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(customCon, &err)
			return
		}

		var sourceProto uint16
		err = binary.Read(conn, binary.LittleEndian, &sourceProto)
		if connectionDown(&err) {
			c.handleConnectionDown(customCon, &err)
			return
		}
		var destProto uint16
		err = binary.Read(conn, binary.LittleEndian, &destProto)
		if connectionDown(&err) {
			c.handleConnectionDown(customCon, &err)
			return
		}
		var dataLen uint32
		err = binary.Read(conn, binary.LittleEndian, &dataLen)
		if connectionDown(&err) {
			c.handleConnectionDown(customCon, &err)
			return
		}
		buffer := make([]byte, dataLen)
		var readData int
		readData, err = conn.Read(buffer)
		if connectionDown(&err) || err != nil {
			c.handleConnectionDown(customCon, &err)
			return
		}
		if readData != int(dataLen) {
			str := fmt.Sprintf("EXPECTED TO READ %d; BUT ONLY READ %d. CONN CLOSED", dataLen, readData)
			err = errors.New(str)
			fmt.Println(str)
			c.handleConnectionDown(customCon, &err)
			return
		}
		netEvent := gobabelUtils.MESSAGE_RECEIVED
		if listenAddress == nil {
			tcpAddr, err := net.ResolveTCPAddr("tcp", string(buffer))
			if err == nil {
				log.Println("THE CLIENT REMOTE ADDRESS IS: ", tcpAddr.String())
				c.auxAddCon(customCon, conn, tcpAddr)
				listenAddress = tcpAddr
				netEvent = gobabelUtils.CONNECTION_UP
			} else {
				_ = conn.Close()
				log.Fatal("FAILED TO PARSE RECEIVED ADDRESS OF THE SERVER:", err)
				return
			}
			buffer = nil
		}
		msgEvent := NewNetworkEvent(customCon, buffer, gobabelUtils.APP_PROTO_ID(sourceProto), gobabelUtils.APP_PROTO_ID(destProto), netEvent, gobabelUtils.MessageHandlerID(msgHandlerId))
		c.protoListener.DeliverEvent(msgEvent)
	}
}

func (c *TCPChannel) SendAppData(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) (int, error) {
	return c.sendMessage(ipAddress, source, destProto, msg, gobabelUtils.APP_MSG, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
}
func (c *TCPChannel) SendAppData2(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg NetworkMessage, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	customWriter := NewCustomWriter3(binary.LittleEndian)
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	writeHeaders(customWriter, source, destProto, gobabelUtils.APP_MSG, msgHandlerId)
	msg.SerializeData(customWriter)
	return c.auxWriteToNetwork(hostAddress, customWriter)
}

func (c *TCPChannel) IsConnected(address string) bool {
	return c.connections[address] != nil
}

func auxWriteToNetworkChild(connection *CustomConnection, writer *CustomWriter) (int, error) {
	dataSize := uint32(writer.Len() - HeaderSize)
	aux := writer.OffSet()
	writer.SetOffSet(HeaderSize - 4)
	_, _ = writer.WriteUInt32(dataSize)
	writer.SetOffSet(aux)
	//println("SENT DATA IS ", base32.HexEncoding.EncodeToString(writer.Data()))
	return connection.conn.Write(writer.Data())
}

func (c *TCPChannel) auxWriteToNetwork(address string, writer *CustomWriter) (int, error) {
	customCon := c.connections[address]
	written := -1
	var err error
	if customCon != nil {
		return auxWriteToNetworkChild(customCon, writer)
		//written to logger service
	} else {
		err = errors.New(fmt.Sprintf("%s IS NOT CONNECTED", address))
	}
	return written, err
}

// TODO use a golang routine just to send messages
// write order: msgType(uint8), sourceProto(uint16), destProto(uint16),dataLength(uint32),appData
func (c *TCPChannel) sendMessage(hostAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, msgType gobabelUtils.MSG_TYPE, msgHandlerId gobabelUtils.MessageHandlerID) (int, error) {
	//I DONT LIKE IT:
	// using protoBuf to binary a struct with the data, msgType, ??
	//TODO ASSERT THAT THE BUFFER IS THE SAME AFTER BEING PASSED TO HEADER AND THEN TO THE SERIALIZEER
	customWriter := NewCustomWriter(HeaderSize+len(msg), binary.LittleEndian)
	writeHeaders(customWriter, source, destProto, msgType, msgHandlerId)
	customWriter.Write(msg)
	return c.auxWriteToNetwork(hostAddress, customWriter)
}
func (t *TCPChannel) handleConnectionDown(customCon *CustomConnection, err *error) {
	//_ = customCon.conn.Close()
	t.onDisconnected(customCon.connectionKey)
	msgEvent := NewNetworkEvent(customCon, nil, gobabelUtils.ALL_PROTO_ID, gobabelUtils.ALL_PROTO_ID, gobabelUtils.CONNECTION_DOWN, gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID)
	t.protoListener.DeliverEvent(msgEvent)
}
func connectionDown(err *error) bool {
	return errors.Is(*err, io.EOF) || errors.Is(*err, syscall.EPIPE)
}

/******************* ******************* ******************* ******************* *******************/

type NetworkEvent struct {
	customConn       *CustomConnection
	Data             []byte
	SourceProto      gobabelUtils.APP_PROTO_ID
	DestProto        gobabelUtils.APP_PROTO_ID
	NET_EVENT        gobabelUtils.NET_EVENT
	MessageHandlerID gobabelUtils.MessageHandlerID
}

func NewNetworkEvent(customCon *CustomConnection, data []byte, sourceProto, destProto gobabelUtils.APP_PROTO_ID, eventType gobabelUtils.NET_EVENT, messageHandlerID gobabelUtils.MessageHandlerID) *NetworkEvent {
	return &NetworkEvent{
		customConn:       customCon,
		Data:             data,
		SourceProto:      sourceProto,
		DestProto:        destProto,
		NET_EVENT:        eventType,
		MessageHandlerID: messageHandlerID,
	}
}
