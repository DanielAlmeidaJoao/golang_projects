package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	quitch := make(chan struct{})

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
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case t := <-ticker.C:
				//fmt.Println("Tick at", t)
				data = fmt.Sprintf("Tick at %s", t.String())
				_, err = io.Copy(conn, strings.NewReader(data))
				if err != nil {
					println("FAILED TO SEND MESSAGE ON TIME TICK")
				}
			}
		}
	}()

	for {
		//Read the echoed data back from the server
		buf := make([]byte, len(data))
		if _, err := io.ReadFull(conn, buf); err != nil {
			fmt.Println("Error reading data:", err)
			return
		}

		fmt.Println(string(buf))
	}

	<-quitch
}
