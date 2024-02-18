package main

import (
	"fmt"
	"net"
)

func handleConnectin(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from", conn.RemoteAddr())

	for {
		//io.ReadAll()
		//io.Copy is blocking
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading data:", err)
			return
		}
		buf = buf[:n]
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("Error echoing data:", err)
			return
		}
		println("Echoed number of characters: ", string(buf))
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
