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
	minCap := c.offset + len(p)
	if minCap >= cap(c.data) {
		aux := c.data
		c.data = make([]byte, minCap+(minCap*1/3))
		copy(c.data, aux)
	}
	for i := range p {
		c.data[c.offset] = p[i]
		c.offset++
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
	if err != nil {
		return 0, err
	}
	return c.offset - old, nil
}
func (c *CustomWriter) SetOffSet(offSet int) int {
	if offSet < cap(c.data) {
		c.offset = offSet
	}

	return c.offset
}

func (c *CustomWriter) OffSet() int {
	return c.offset
}

func (c *CustomWriter) Len() int {
	return len(c.data)
}

func (c *CustomWriter) Data() []byte {
	return c.data[:c.offset]
}
