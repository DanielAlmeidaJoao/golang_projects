package protocolLIstenerLogics

import (
	gobabelUtils "gobabel/commons"

	"net"
)

type ProtoInterface interface {
	ProtocolUniqueId() gobabelUtils.APP_PROTO_ID
	OnStart(channelInterface ChannelInterface)
	OnMessageArrival(from *net.Addr, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte, channelInterface ChannelInterface)
	ConnectionUp(from *net.Addr, channelInterface ChannelInterface)
	ConnectionDown(from *net.Addr, channelInterface ChannelInterface)
}

type ProtocolManager struct {
	id            gobabelUtils.APP_PROTO_ID
	eventQueue    <-chan *gobabelUtils.NetworkEvent
	handlerID     gobabelUtils.MessageHandlerID
	eventHandlers map[gobabelUtils.MessageHandlerID]any
	channel       ChannelInterface
}

func (p *ProtocolManager) connectionUp(from *net.Addr) {

}
func (p *ProtocolManager) connectionDown(from *net.Addr) {

}
func (p *ProtocolManager) sendMessage(ipAddress string, source, destProto gobabelUtils.APP_PROTO_ID, msg []byte) {

}
func (p *ProtocolManager) deliveryEvent(networkEvent gobabelUtils.NetworkEvent) {

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
