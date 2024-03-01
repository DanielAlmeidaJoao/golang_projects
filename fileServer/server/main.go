package main

//nc -vz localhost 3000
//we need to link all files, go run *.main or go run .
import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func abort(err error) bool {
	if err == nil {
		return false
	} else {
		fmt.Println("error: ", err)
		os.Exit(0)
		return true

	}
}
func connectionDown(err error) bool {
	return errors.Is(err, io.EOF)
}
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2024) //dangerous

	newFileName := "newFile_copy2"
	file, err := os.Create(newFileName)
	defer closeFile(file)
	if err != nil {
		println("failed to open a connection: ", err)
		return
	}
	for {
		readBytes, err := conn.Read(buffer)
		if connectionDown(err) {
			println("CLIENT CONNECTION IS DOWN!")
		} else if readBytes > 0 {
			aux := buffer[:readBytes]
			println("RECEIVED BYTES: ", readBytes)
			_, err := file.Write(aux)
			println("GOING TO SLEEP ...")
			time.Sleep(time.Minute * 1)
			if err != nil {
				println("FAILED TO WRITE: ", err)
			}
		}

	}
}
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		abort(err)
	}
}
func sendFile(conn net.Conn) {
	file, err := os.Open("/Users/danieljoao/Documents/bla copy.json")
	defer closeFile(file)
	if abort(err) {
		os.Exit(0)
	}

	for {
		buffer := make([]byte, 2024)

		size, err := file.Read(buffer)

		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		abort(err)
		if size > 0 {
			n, err := conn.Write(buffer[:size])
			if err != nil {
				println("CLIENT DOWN. FAILED TO WRITE")
			}
			println("NUMBER OF BYTES WRITTEN: ", n)
		} else {
			println("NUMBER OF BYTES WRITTEN: ", size)
			break
		}
	}
	println("FIINISHED SENDING THE WHOLE FILE")
}
func main() {
	listener, err := net.Listen("tcp", ":3000")
	println("SERVER STARTED. WAITING FOR CONNECTIONS ...")
	if abort(err) {
		return
	}
	for {
		conn, err := listener.Accept()
		println("SERVER ACCEPTED A CONNECTION:", conn.RemoteAddr())
		if abort(err) {
			return
		}

		go handleConnection(conn)
		go sendFile(conn)

	}

	fmt.Println("END ...")
}
