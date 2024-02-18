package main

import (
	"fmt"
	"io"
	"net"
)

func handleConnectin(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr())

	for {
		//io.ReadAll()
		//io.Copy is blocking
		written, err := io.Copy(conn, conn)

		if err != nil {
			fmt.Println("Error echoing data:", err)
			return
		}
		println("Echoed number of characters:", written)

	}
}

func main() {
	//create a TCP listener on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Listening on :8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		//handle the connection in a separate goroutine
		handleConnectin(conn)
	}

}
