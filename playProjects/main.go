package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func sayBlalalal() {
	fmt.Println("HELLO, THE FUNCTION WAS TRIGGERED")
}
func testAfterFuncTime(data interface{}) {

}

func testDefer() {
	ola := false

	defer func(k bool) {
		println(k)
	}(ola)

	ola = true
}
func testIntConversion() {
	buf := new(bytes.Buffer)
	var num int = 1
	err := binary.Write(buf, binary.LittleEndian, num)

	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	var dest int
	err = binary.Read(buf, binary.LittleEndian, &dest)
	println(err)
	println("PRIGINAL IS:", dest)
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

	n := []int{1, 2, 3, 4, 5, 6}
	n2 := make([]int, 10)
	copy(n2, n)

	fmt.Println("N2: IS :", n2)
	copy(n2, n)
	fmt.Println("N23: IS :", n2)

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

func testComplex() {

}

type OLAINTERFACE interface {
	sayOla() string
	setOla(ol string) string
}
type ImplOla struct {
	ola string
}

func (p *ImplOla) setOla(ol string) string {
	p.ola = ol
	return p.ola
}

func (p *ImplOla) sayOla() string {
	return p.ola
}

type NetworkMessageDeserializer[T OLAINTERFACE] interface {
	deserializeData(reader string) *T //make it return the
}

func changeInterfaceValue(olainterface OLAINTERFACE) {
	olainterface.setOla("CARALHOO")
}
func changeInterfaceValue2(olainterface *OLAINTERFACE) {
	(*olainterface).setOla("PORRASS")

}
func testOla() {
	gg := ImplOla{
		ola: "OLAKK",
	}
	hh := &gg
	var jj OLAINTERFACE = hh

	println(hh.sayOla())
	println(gg.sayOla())

	changeInterfaceValue(hh)
	changeInterfaceValue2(&jj)

	println(hh.sayOla())
	println(gg.sayOla())

}
func testSlices() {
	o := []int{0, 0, 0}
	p := []int{1, 2, 3}
	k := append(o, p...)
	for i := range k {
		println(k[i])
	}
}

func testLogTrace() {
	log.Println("ERROR HAPPENED ---")
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testDefer()
	//testBytesConversion()
	//testAppendConversion()
	//testMapLen()
	//testGenerics()
	//testIntConversion()
	//testSlices()
	//testOla()
	//testAfterFuncTime()
	testLogTrace()
}
