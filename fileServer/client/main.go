package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

func abort(err error) bool {
	if err == nil {
		return false
	} else {
		fmt.Println("error: ", err)
		return true
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2024) //dangerous

	for {
		readBytes, err := conn.Read(buffer)
		if abort(err) {
			return
		}
		str := string(buffer[:readBytes])
		println("READ: ", str)
	}
}

func main() {
	totalRead := 0
	conn, err := net.Dial("tcp", ":3000")
	defer conn.Close()
	if abort(err) {
		return
	}
	//text := "TODAY IS A GOOD DAY! I CAN´T BELIEVE IT!!!"
	//conn.Write([]byte(text))

	buffer := make([]byte, 2024) //dangerous

	for {
		readBytes, err := conn.Read(buffer)
		totalRead += readBytes
		if errors.Is(err, io.EOF) {
			println("SERVER IS DOWN")
			return
		}
		abort(err)
		println("CLIENT READ AND SENT: ", readBytes)
		_, err = conn.Write(buffer[:readBytes])
		abort(err)
	}

}
