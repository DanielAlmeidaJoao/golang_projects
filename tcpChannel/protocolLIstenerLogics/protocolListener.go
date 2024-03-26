package protocolLIstenerLogics

import (
	"encoding/binary"
	"fmt"
	gobabelUtils "gobabel/commons"
	"log"
	"net"
	"sync"
	"time"
)

type ConnectionState struct {
	address net.Addr
	state   gobabelUtils.NET_EVENT
	err     error
}
type ProtoListenerInterface interface {
	AddProtocol(protocol ProtoInterface) error
	Start() error
	RegisterNetworkMessageHandler(handlerId gobabelUtils.MessageHandlerID, funcHandler MESSAGE_HANDLER_TYPE) error
	RegisterTimeout(sourceProto gobabelUtils.APP_PROTO_ID, duration time.Duration, data interface{}, funcToExecute gobabelUtils.TimerHandlerFunc) int
	RegisterLocalCommunication(sourceProto, destProto gobabelUtils.APP_PROTO_ID, data interface{}, funcToExecute gobabelUtils.LocalProtoComHandlerFunc) error
	CancelTimer(timerId int) bool
	RegisterPeriodicTimeout(sourceProto gobabelUtils.APP_PROTO_ID, duration time.Duration, data interface{}, funcToExecute gobabelUtils.TimerHandlerFunc) int
}
type protoWrapper struct {
	queue                   chan *gobabelUtils.NetworkEvent
	timeoutChannel          chan int
	localCommunicationQueue chan *gobabelUtils.LocalCommunicationEvent
	proto                   ProtoInterface
}
type CustomPair[F any, S any] struct {
	First  F
	Second S
}
type timerArgs struct {
	protoId       gobabelUtils.APP_PROTO_ID
	data          interface{}
	funcHandler   gobabelUtils.TimerHandlerFunc
	timer         *time.Timer
	periodicTimer *time.Ticker
}
type ProtoListener struct {
	mutex           sync.Mutex
	protocols       map[gobabelUtils.APP_PROTO_ID]*protoWrapper
	channel         ChannelInterface
	messageHandlers map[gobabelUtils.MessageHandlerID]MESSAGE_HANDLER_TYPE
	timerHandlers   map[int]*timerArgs
	timersId        int
	order           binary.ByteOrder
	ConnectionType  gobabelUtils.CONNECTION_TYPE
}

/*********************** CLIENT METHODS ***************************/
func (p *ProtoListener) RegisterNetworkMessageHandler(handlerId gobabelUtils.MessageHandlerID, funcHandler MESSAGE_HANDLER_TYPE) error {
	var err error
	p.mutex.Lock()
	if p.messageHandlers[handlerId] == nil {
		p.messageHandlers[handlerId] = funcHandler
		err = nil
	} else {
		err = gobabelUtils.ELEMENT_EXISTS_ALREADY
	}
	p.mutex.Unlock()
	return err
}

// address string, port int, connectionType gobabelUtils.CONNECTION_TYPE
func NewProtocolListener(address string, port int, connectionType gobabelUtils.CONNECTION_TYPE, order binary.ByteOrder) ProtoListenerInterface {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ch := NewTCPChannel(address, port, connectionType)
	protoL := &ProtoListener{
		protocols:       make(map[gobabelUtils.APP_PROTO_ID]*protoWrapper),
		timerHandlers:   make(map[int]*timerArgs),
		channel:         ch,
		order:           order,
		messageHandlers: make(map[gobabelUtils.MessageHandlerID]MESSAGE_HANDLER_TYPE),
		ConnectionType:  connectionType,
	}
	ch.SetProtoLister(protoL)
	return protoL
}

func (l *ProtoListener) RegisterTimeout(sourceProto gobabelUtils.APP_PROTO_ID, duration time.Duration, data interface{}, funcToExecute gobabelUtils.TimerHandlerFunc) int {
	aux := l.timersId
	t := time.AfterFunc(duration, func() {
		l.protocols[sourceProto].timeoutChannel <- aux
	})
	l.timerHandlers[l.timersId] = &timerArgs{
		protoId:     sourceProto,
		data:        data,
		funcHandler: funcToExecute,
		timer:       t,
	}
	l.timersId++
	return aux
}

// I DONT LIKE IT BECAUSE FOR EACH TIMER I NEED TO HAVE A GOROUTINE
func (l *ProtoListener) RegisterPeriodicTimeout(sourceProto gobabelUtils.APP_PROTO_ID, duration time.Duration, data interface{}, funcToExecute gobabelUtils.TimerHandlerFunc) int {
	ticker := time.NewTicker(duration)
	aux := l.timersId
	go func() {
		for range ticker.C {
			l.protocols[sourceProto].timeoutChannel <- aux
		}
	}()
	l.timerHandlers[l.timersId] = &timerArgs{
		protoId:       sourceProto,
		data:          data,
		funcHandler:   funcToExecute,
		timer:         nil,
		periodicTimer: ticker,
	}
	l.timersId++
	return aux
}

func (l *ProtoListener) RegisterLocalCommunication(sourceProto, destProto gobabelUtils.APP_PROTO_ID, data interface{}, funcToExecute gobabelUtils.LocalProtoComHandlerFunc) error {
	proto := l.protocols[destProto]
	if proto == nil {
		return gobabelUtils.UNKNOWN_PROTOCOL
	}
	proto.localCommunicationQueue <- gobabelUtils.NewLocalCommunicationEvent(sourceProto, destProto, data, funcToExecute)
	return nil
}

func (l *ProtoListener) CancelTimer(timerId int) bool {
	args := l.timerHandlers[timerId]
	if args != nil {
		if args.timer == nil {
			args.periodicTimer.Stop()
		} else {
			args.timer.Stop()
		}
	}
	return true
}

func (l *ProtoListener) AddProtocol(protocol ProtoInterface) error {
	//TODO make the constants dynamic
	if l.protocols[(protocol).ProtocolUniqueId()] != nil {
		return gobabelUtils.PROTOCOL_EXIST_ALREADY
	}
	l.protocols[(protocol).ProtocolUniqueId()] = &protoWrapper{
		queue:                   make(chan *gobabelUtils.NetworkEvent, 50),
		proto:                   protocol,
		timeoutChannel:          make(chan int, 100),
		localCommunicationQueue: make(chan *gobabelUtils.LocalCommunicationEvent, 20),
	}
	return nil
}

// TODO should all protocols receive connection up event ??
func (l *ProtoListener) Start() error {
	if len(l.protocols) == 0 {
		log.Fatal(gobabelUtils.NO_PROTOCOLS_TO_RUN)
		return gobabelUtils.NO_PROTOCOLS_TO_RUN
	}
	for _, protoWrapper := range l.protocols {
		protoWrapper := protoWrapper
		go func() {
			proto := protoWrapper.proto
			proto.OnStart(l.channel)
			log.Printf("PROTOCOL <%d> STARTED LISTENING TO EVENTS...\n", proto.ProtocolUniqueId())
			for {
				select {
				case networkEvent := <-protoWrapper.queue:
					log.Printf("PROTOCOL <%d> RECEIVED AN EVENT. EVENT TYPE <%d>\n", proto.ProtocolUniqueId(), networkEvent.NET_EVENT)
					switch networkEvent.NET_EVENT {
					case gobabelUtils.CONNECTION_UP:
						proto.ConnectionUp(&networkEvent.From, l.channel)
					case gobabelUtils.CONNECTION_DOWN:
						proto.ConnectionDown(&networkEvent.From, l.channel)
					case gobabelUtils.MESSAGE_RECEIVED:
						if gobabelUtils.NO_NETWORK_MESSAGE_HANDLER_ID == networkEvent.MessageHandlerID {
							proto.OnMessageArrival(&networkEvent.From, networkEvent.SourceProto, networkEvent.DestProto, networkEvent.Data, l.channel)
						} else {
							messageHandler := l.messageHandlers[networkEvent.MessageHandlerID]
							messageHandler(networkEvent.From.String(), networkEvent.SourceProto, NewCustomReader(networkEvent.Data, l.order))
						}
					default:
						(l.channel).CloseConnection(networkEvent.From.String())
						log.Fatal(fmt.Sprintf("RECEIVED AN EVENT NOT PART OF THE PROTOCOL. CONNECTION CLOSED %s", networkEvent.From.String()))
					}
				case timerId := <-protoWrapper.timeoutChannel:
					args := l.timerHandlers[timerId]
					if args != nil {
						if args.timer != nil {
							// it is a timeout, otherwise it is a periodic timer
							delete(l.timerHandlers, timerId)
						}
						args.funcHandler(args.protoId, args.data)
					}
				case localEventCom := <-protoWrapper.localCommunicationQueue:
					localEventCom.ExecuteFunc()
				}
			}
		}()

	}

	return nil

}

/*********************** CLIENT METHODS ***************************/

// TODO should all protocols receive connection up event ??
func (l *ProtoListener) DeliverEvent(event *gobabelUtils.NetworkEvent) {
	if event.DestProto == gobabelUtils.ALL_PROTO_ID {
		log.Default().Println("GOING TO DELIVER AN EVENT TO ALL PROTOCOLS")
		for _, proto := range l.protocols {
			proto.queue <- event
		}
	} else {
		protocol := l.protocols[event.SourceProto]
		if protocol == nil {
			log.Fatalln("RECEIVED EVENT FOR A NON EXISTENT PROTOCOL!")
		} else {
			log.Default().Println("GOING TO DELIVER AN EVENT TO THE PROTOCOL:", event.SourceProto)
			protocol.queue <- event
		}
	}
}
