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
func main() {
	//testDefer()

	buf := new(bytes.Buffer)
	fmt.Println("DEFAULT SIZE: ", buf.Cap(), buf.Available(), buf.AvailableBuffer())
	var num uint8 = 1
	err := binary.Write(buf, binary.LittleEndian, num)
	fmt.Println("DEFAULT SIZE: ", buf.Cap(), buf.Available(), buf.AvailableBuffer())

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("gggg %x h \n", buf.Bytes())
	fmt.Printf("lennnn %d h \n", len(buf.Bytes()))

	buf.WriteString("hello")
	println()
	fmt.Println(buf.Available(), "olll")
	fmt.Println(buf.Cap())
	var dest uint8
	err = binary.Read(buf, binary.LittleEndian, &dest)
	println(err)
	println("PRIGINAL IS:", dest)
}
