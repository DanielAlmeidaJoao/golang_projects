package protocolLIstenerLogics

import (
	"fmt"
	gobabelUtils "gobabel/commons"
	"log"
	"net"
)

//address string, port int, connectionType gobabelUtils.CONNECTION_TYPE
func NewProtocolListener(address string, port int, connectionType gobabelUtils.CONNECTION_TYPE) ProtoListenerInterface {
	ch := NewTCPChannel(address, port, connectionType)
	protoL := &ProtoListener{
		protocols: make(map[gobabelUtils.APP_PROTO_ID]*protoWrapper),
		channel:   ch,
	}
	ch.SetProtoLister(protoL)
	return protoL
}

type ConnectionState struct {
	address net.Addr
	state   gobabelUtils.NET_EVENT
	err     error
}
type ProtoListenerInterface interface {
	AddProtocol(protocol ProtoInterface) error
	Start() error
}
type protoWrapper struct {
	queue chan *gobabelUtils.NetworkEvent
	proto ProtoInterface
}
type ProtoListener struct {
	protocols map[gobabelUtils.APP_PROTO_ID]*protoWrapper
	channel   ChannelInterface
}

func (l *ProtoListener) AddProtocol(protocol ProtoInterface) error {
	if l.protocols[(protocol).ProtocolUniqueId()] != nil {
		return gobabelUtils.PROTOCOL_EXIST_ALREADY
	}
	l.protocols[(protocol).ProtocolUniqueId()] = &protoWrapper{
		queue: make(chan *gobabelUtils.NetworkEvent, 20),
		proto: protocol,
	}
	return nil
}

//TODO should all protocols receive connection up event ??
func (l *ProtoListener) Start() error {
	if len(l.protocols) == 0 {
		log.Fatal(gobabelUtils.NO_PROTOCOLS_TO_RUN)
		return gobabelUtils.NO_PROTOCOLS_TO_RUN
	}
	for _, protoWrapper := range l.protocols {
		go func() {
			proto := protoWrapper.proto
			proto.OnStart(l.channel)
			log.Printf("PROTOCOL <%d> STARTED LISTENING TO EVENTS...\n", proto.ProtocolUniqueId())
			for {
				networkEvent := <-protoWrapper.queue
				log.Printf("PROTOCOL <%d> RECEIVED AN EVENT. EVENT TYPE <%d>\n", proto.ProtocolUniqueId(), networkEvent.NET_EVENT)
				switch networkEvent.NET_EVENT {
				case gobabelUtils.CONNECTION_UP:
					proto.ConnectionUp(&networkEvent.From, l.channel)
				case gobabelUtils.CONNECTION_DOWN:
					proto.ConnectionDown(&networkEvent.From, l.channel)
				case gobabelUtils.MESSAGE_RECEIVED:
					proto.OnMessageArrival(&networkEvent.From, networkEvent.SourceProto, networkEvent.DestProto, networkEvent.Data, l.channel)
				default:
					(l.channel).CloseConnection(networkEvent.From.String())
					log.Fatal(fmt.Sprintf("RECEIVED AN EVENT NOT PART OF THE PROTOCOL. CONNECTION CLOSED %s", networkEvent.From.String()))
				}
			}
		}()

	}

	return nil

}

//TODO should all protocols receive connection up event ??
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
