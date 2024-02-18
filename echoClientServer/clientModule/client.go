package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	fmt.Println("Connected to server")

	//Send some data to the server
	data := "Hello, server!"
	_, err = io.Copy(conn, strings.NewReader(data))
	fmt.Println("READ FROM SERVER")
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	//Read the echoed data back from the server
	buf := make([]byte, len(data))
	if _, err := io.ReadFull(conn, buf); err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}
	fmt.Printf("Received from server:", string(buf))
}
