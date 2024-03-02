package main

//nc -vz localhost 3000
//we need to link all files, go run *.main or go run .
import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
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
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE)
}
func handleConnection(conn net.Conn, endChan chan<- int) {

	buffer := make([]byte, 1024*64) //dangerous
	newFileName := time.Now().String() + conn.LocalAddr().String()
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
			return
		}
		if readBytes > 0 {
			println("Read bytes: ", readBytes)
			endChan <- readBytes
			aux := buffer[:readBytes]
			_, err = file.Write(aux)
			if connectionDown(err) {
				return
			}
			println("RECEIVED BYTES: ", readBytes)
		}

	}
}
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		abort(err)
	}
}
func sendFile(conn net.Conn, endChan chan<- int) {
	//dir, _ := os.Getwd()
	//path := fmt.Sprintf("%s/testFiles/newFile_copy2", dir)

	path := "/home/tsunami/Downloads/plane/movie.mp4"
	file, err := os.Open(path)
	defer closeFile(file)
	if abort(err) {
		os.Exit(0)
	}
	totalSent := 0
	for {
		buffer := make([]byte, 1024*128)

		size, err := file.Read(buffer)
		totalSent += size
		if errors.Is(err, io.EOF) {
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
	endChan <- totalSent
	println("FIINISHED SENDING THE WHOLE FILE. TOTAL SENT: ", totalSent)
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
		sendChan := make(chan int)
		receiveChan := make(chan int)

		go handleConnection(conn, receiveChan)
		go sendFile(conn, sendChan)
		go func() {
			sent, received := 0, 0
			for {
				select {
				case sent = <-sendChan:
				case x := <-receiveChan:
					received += x
				}
				if sent > 0 && sent == received {
					println("SERVER RECEIVED THE SAME AMOUNT OF BYTES SENT TO THE CLIENT: ", sent)
					conn.Close()
					return
				}
			}
		}()
	}

	fmt.Println("END ...")
}
