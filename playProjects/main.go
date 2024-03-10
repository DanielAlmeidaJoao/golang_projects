package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func testDefer() {
	ola := false

	defer func(k bool) {
		println(k)
	}(ola)

	ola = true
}
func testBytesConversion() {
	buf := new(bytes.Buffer)
	var num uint8 = 1
	err := binary.Write(buf, binary.LittleEndian, num)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("gggg %x h \n", buf.Bytes())
	fmt.Printf("lennnn %d h \n", len(buf.Bytes()))

	buf.WriteString("hello")
	println()
	fmt.Println(buf.Cap())
	var dest uint8
	err = binary.Read(buf, binary.LittleEndian, &dest)
	println(err)
	println("PRIGINAL IS:", dest)
}
func testAppendConversion() {
	numbers := make([]int, 10)
	fmt.Println(numbers)
	numbers = append(numbers, 123)
	numbers = append(numbers, 99)

	fmt.Println(numbers)

}
func testMapLen() {
	ma := make(map[int]string, 10)
	fmt.Println(len(ma))

	ma[12] = "OLA"
	ma[12] = "OLA2"
	ma[13] = "PPOLA"
	fmt.Println(ma)
}

type EventHolder[T any] struct {
	eventId int
	data    T
}
type Message struct {
	name   string
	age    int
	weight int
}

func testGenerics() {
	event := EventHolder[Message]{
		eventId: 100,
		data: Message{
			age:    12,
			name:   "DANIEL",
			weight: 123,
		},
	}

	fmt.Println(event.eventId, event.data.name, event.data.age, event.data.weight)
}
func main() {
	//testDefer()
	//testBytesConversion()
	//testAppendConversion()
	//testMapLen()
	testGenerics()
}
