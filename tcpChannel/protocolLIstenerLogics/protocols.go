package protocolLIstenerLogics

type Protocol struct {
	id            int16
	eventQueue    chan *NetworkEvent
	handlerID     eventHandlerId
	eventHandlers map[eventHandlerId]any
}

func (p *Protocol) deliveryEvent() {

}
func (p *Protocol) nextEventId() eventHandlerId {
	p.handlerID++
	return p.handlerID
}
func (p *Protocol) start() {
	for {
		event := <-p.eventQueue
		networkEventHandler := p.eventHandlers[event.EventId]
		if networkEventHandler != nil {
			// TODO networkEventHandler(event)
		}
	}
}
