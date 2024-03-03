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
	var num uint16 = 1234
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
	var dest uint16
	err = binary.Read(buf, binary.LittleEndian, &dest)
	println(err)
	println("PRIGINAL IS:", dest)
}
