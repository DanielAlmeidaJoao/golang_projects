package protocolLIstenerLogics

import (
	"encoding/binary"
	gobabelUtils "gobabel/commons"
)

type NetworkMessage interface {
	SerializeData(writer *CustomWriter)
}

// type MESSAGE_HANDLER_TYPE func(from string, protoSource APP_PROTO_ID, data []byte)
type MESSAGE_HANDLER_TYPE func(from string, protoSource gobabelUtils.APP_PROTO_ID, data NetworkMessage)

type CustomWriter struct {
	data   []byte
	offset int
	order  binary.ByteOrder
}

func NewCustomWriter(capacity int, order binary.ByteOrder) *CustomWriter {
	return &CustomWriter{
		data:   make([]byte, capacity),
		offset: 0,
		order:  order,
	}
}
func NewCustomWriter2(data []byte, offset int, order binary.ByteOrder) *CustomWriter {
	return &CustomWriter{
		data:   data,
		offset: offset,
		order:  order,
	}
}
func NewCustomWriter3(order binary.ByteOrder) *CustomWriter {
	return NewCustomWriter(64, order)
}
func (c *CustomWriter) Write(p []byte) (n int, err error) {
	if len(p) > (cap(c.data) - len(c.data)) {
		c.data = append(c.data, p...)
		c.offset = len(c.data)
	} else {
		for i := range p {
			c.data[c.offset] = p[i]
			c.offset++
		}
	}
	return len(p), nil
}

func (c *CustomWriter) WriteNumber(p any) (n int, err error) {
	old := c.offset
	er := binary.Write(c, c.order, p)
	return c.offset - old, er
}
func (c *CustomWriter) WriteString(p string) (n int, err error) {
	old := c.offset
	c.Write([]byte(p))
	return c.offset - old, nil
}
func (c *CustomWriter) SetOffSet(offSet int) int {
	if offSet < cap(c.data) {
		c.offset = offSet
	}

	return c.offset
}

func (c *CustomWriter) Len() int {
	return len(c.data)
}
