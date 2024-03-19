package protocolLIstenerLogics

import "encoding/binary"

type MESSAGE_DESERIALIZER_TYPE func(reader *CustomReader) NetworkMessage

/**
type NetworkMessageDeserializer interface {
	DeserializeData(reader *CustomReader) NetworkMessage //make it return the
} **/

type CustomReader struct {
	data   []byte
	offset int
	order  binary.ByteOrder
}

func (c *CustomReader) Read(p []byte) (n int, err error) {
	c.offset += copy(p, c.data[c.offset:])
	return c.offset, nil
}

/*
*
Writes any number like int32, uint8, ... to the buffer. Numbers like int or float are not accepted
*/
func (b *CustomReader) ReadAnyNumberWithSizeOnTheName(dest any) error {
	err := binary.Read(b, b.order, dest)
	return err
}
func (b *CustomReader) ReadUint32() (uint32, error) {
	var value uint32
	err := binary.Read(b, b.order, value)
	b.offset += 4
	return value, err
}
func (b *CustomReader) ReadAll() []byte {
	return b.data[b.offset:]
}
func (b *CustomReader) ReadString(n int) (string, error) {
	//TODO check for error
	aux := b.data[b.offset : b.offset+n]
	b.offset += len(aux)
	return string(aux), nil
}
func NewCustomReader(data []byte, order binary.ByteOrder) *CustomReader {
	return &CustomReader{
		data:   data,
		offset: 0,
		order:  order,
	}
}
