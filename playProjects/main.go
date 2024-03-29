package main

import (
	"bytes"
	list2 "container/list"
	"encoding/binary"
	"fmt"
	"log"
	"net"
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
	fmt.Println(len(ma))
	delete(ma, 12)
	fmt.Println(len(ma))

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
func testList() {
	list := list2.New()
	list.PushBack(1)
	list.PushBack(2)
	list.PushBack(3)
	list.PushBack(4)
	list.PushBack(5)
	list.PushBack(6)
	elem := list.Front()
	fmt.Println(elem.Value)

	i := 0
	next := list.Front()
	var del *list2.Element
	for i < list.Len() {
		fmt.Println("ELEMENTS OF THE LIST: ", next.Value)
		if next.Value == 2 {
			del = next
		}
		next = next.Next()
		i++
	}
	fmt.Println("REMOVED ELEMENT: ", list.Remove(del))
	next = list.Front()
	i = 0
	for i < list.Len() {
		fmt.Println("ELEMENTS OF THE LIST: ", next.Value)
		next = next.Next()
		i++
	}
}

func testAddressesDifs() {
	var ad map[net.Addr]int = make(map[net.Addr]int)

	fmt.Println(ad)
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	fmt.Println(tcpAddr, err)
	tcpAddr2, err2 := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	fmt.Println(tcpAddr2, err2)
	ad[tcpAddr] = 1232
	fmt.Println(tcpAddr == tcpAddr2)
	fmt.Println(tcpAddr.String() == tcpAddr2.String())
	fmt.Println(tcpAddr2.String())
	fmt.Println(ad[tcpAddr2])

}
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//testDefer()
	//testBytesConversion()
	//testAppendConversion()
	testMapLen()
	//testGenerics()
	//testIntConversion()
	//testSlices()
	//testOla()
	//testAfterFuncTime()
	//testLogTrace()
	//testList()
	//testAddressesDifs()
}
