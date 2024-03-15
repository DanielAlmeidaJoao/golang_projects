package protocolLIstenerLogics

//BIG TODO MAKE THIS GENERIC LIKE IN THE GENERICS YOUTUBE VIDEO
import (
	"bytes"
	"encoding/binary"
	"golang.org/x/exp/constraints"
	_ "golang.org/x/exp/constraints"
	"io"
)

type CustomData interface {
	constraints.Complex
}
type TestGenNumb[T constraints.Unsigned] struct {
	num T
}

type BinaryWriter struct {
	buf       *bytes.Buffer
	byteOrder binary.ByteOrder
}

type CustomWriter struct {
	data   []byte
	offset int
}

func NewCustomWriter(capacity int) *CustomWriter {
	return &CustomWriter{
		data:   make([]byte, capacity),
		offset: 0,
	}
}
func NewCustomWriter2(data []byte, offset int) *CustomWriter {
	return &CustomWriter{
		data:   data,
		offset: offset,
	}
}
func (c *CustomWriter) Write(p []byte) (n int, err error) {
	if cap(p) > len(c.data) {
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

func (c *CustomWriter) SetOffSet(offSet int) int {
	if offSet < cap(c.data) {
		c.offset = offSet
	}
	return c.offset
}

func (b *BinaryWriter) WriteInt32(value int32) error {
	return binary.Write(b.buf, b.byteOrder, value)
}

/*
*
Writes any number like int32, uint8, ... to the buffer. Numbers like int or float are not accepted
*/
func (b *BinaryWriter) WriteAnyNumberWithSizeOnTheName(value any) error {
	return binary.Write(b.buf, b.byteOrder, value)
}

func (b *BinaryWriter) WriteString(value string) error {
	n, err := b.buf.WriteString(value)
	if n != len(value) {
		return err
	}
	return err
}
func (b *BinaryWriter) writeBinary(value []byte) (int, error) {
	return b.buf.Write(value)
}
func (b *BinaryWriter) writeTo(dest io.Writer) (int64, error) {
	return b.buf.WriteTo(dest)
}

/***************************** BINARY READER *******************************/

type BinaryReader struct {
	buf       io.Reader
	byteOrder binary.ByteOrder
}

func NewBinaryWriter(order binary.ByteOrder) *BinaryWriter {
	return &BinaryWriter{buf: new(bytes.Buffer),
		byteOrder: order}
}
func NewBinaryWriter3(order binary.ByteOrder) *BinaryWriter {
	return &BinaryWriter{buf: new(bytes.Buffer),
		byteOrder: order}
}

func NewBinaryWriter2(buf *bytes.Buffer, order binary.ByteOrder) *BinaryWriter {
	return &BinaryWriter{buf: buf,
		byteOrder: order}
}

func NewBinaryReader(order binary.ByteOrder, reader io.Reader) *BinaryReader {
	return &BinaryReader{buf: reader,
		byteOrder: order}
}

func (b *BinaryReader) ReadInt32() (int32, error) {
	var value int32
	err := binary.Read(b.buf, b.byteOrder, value)
	return value, err
}

/*
*
Writes any number like int32, uint8, ... to the buffer. Numbers like int or float are not accepted
*/
func (b *BinaryReader) ReadAnyNumberWithSizeOnTheName() (any, error) {
	var value any
	err := binary.Read(b.buf, b.byteOrder, value)
	return value, err
}

/************************* NETWORK MESSAGES **************************/

type NetworkMessage interface {
	serializeData(writer *BinaryWriter)
}
type NetworkMessageDeserializer[T NetworkMessage] interface {
	deserializeData(reader *BinaryReader) *T //make it return the
}

//ALWAYS SEND A DEFAULT MESSAGE HANDLERID, IN CASE THE DEVELOPER DOES NOT WANT OT USE THE NETWORKMESSAGE STRUCT
