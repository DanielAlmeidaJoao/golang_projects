package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
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
func connectionDown(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)
}
func main() {
	totalRead := 0
	conn, err := net.Dial("tcp", ":3000")
	defer conn.Close()
	if abort(err) {
		return
	}

	buffer := make([]byte, 1024*64) //dangerous
	written := -1
	for {
		readBytes, err := conn.Read(buffer)
		totalRead += readBytes
		if connectionDown(err) {
			println("REMOTE CONNECTION IS DOWN!")
			return
		}
		println("DATA READ: ", totalRead)
		written, err = conn.Write(buffer[:readBytes])
		if connectionDown(err) {
			return
		}
		println("DATA WRITTEN: ", written)
	}

}
