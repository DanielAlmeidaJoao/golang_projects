package protocolLIstenerLogics

import (
	gobabelUtils "github.com/DanielAlmeidaJoao/golang_projects/tree/main/tcpChannel/gobabel/commons"
)

type ProtoInterface interface {
	ProtocolUniqueId() gobabelUtils.APP_PROTO_ID
	OnStart(channelInterface ChannelInterface)
	OnMessageArrival(customCon *CustomConnection, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface ChannelInterface)
	ConnectionUp(customCon *CustomConnection, channelInterface ChannelInterface)
	ConnectionDown(customCon *CustomConnection, channelInterface ChannelInterface)
}

type ProtocolManager struct {
	id            gobabelUtils.APP_PROTO_ID
	eventQueue    <-chan *NetworkEvent
	handlerID     gobabelUtils.MessageHandlerID
	eventHandlers map[gobabelUtils.MessageHandlerID]any
	channel       ChannelInterface
}

func (p *ProtocolManager) connectionUp(networkEvent *NetworkEvent) {

}
func (p *ProtocolManager) connectionDown(networkEvent *NetworkEvent) {

}
func (p *ProtocolManager) sendMessage(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) {

}
func (p *ProtocolManager) deliveryEvent(networkEvent *NetworkEvent) {

}

func (p *ProtocolManager) nextEventId() gobabelUtils.MessageHandlerID {
	p.handlerID++
	return p.handlerID
}
func (p *ProtocolManager) start() {
	for {
		event := <-p.eventQueue
		networkEventHandler := p.eventHandlers[event.MessageHandlerID]
		if networkEventHandler != nil {
			// TODO networkEventHandler(event)
		}
	}
}
