package protocolLIstenerLogics

import "fileServer/protocolLIstenerLogics/commons"

type Protocol struct {
	id            int16
	eventQueue    chan NetworkEvent
	handlerID     commons.EventID
	eventHandlers map[commons.EventID]any
}

func (p *Protocol) deliveryEvent() {

}
func (p *Protocol) nextEventId() commons.EventID {
	p.handlerID++
	return p.handlerID
}
func (p *Protocol) start() {
	for {
		event := <-p.eventQueue
		networkEventHandler := p.eventHandlers[event.eventId]
		if networkEventHandler != nil {
			// TODO networkEventHandler(event)
		}
	}
}
