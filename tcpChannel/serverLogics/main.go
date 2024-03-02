package main

import (
	"errors"
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"io"
	"net"
	"os"
	"syscall"
)

//go get "github.com/orcaman/concurrent-map"
//nc -vz localhost 3000
//we need to link all files, go run *.main or go run .

const MAX_BYTES_TO_READ = 64 * 1024

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

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, MAX_BYTES_TO_READ)
		_, err := conn.Read(buffer)
		if connectionDown(err) {
			println("CLIENT CONNECTION IS DOWN!")
			return
		}
	}
}
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		abort(err)
	}
}

type Server struct {
	connections cmap.ConcurrentMap
	address     string
	port        int
}

func newServer(address string, port int) *Server {
	m := cmap.New()
	println(m)
	return &Server{
		connections: cmap.New(),
		address:     address,
		port:        port,
	}
}
func (s *Server) start() {
	m := cmap.New()
	println(m)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	abort(err)
	defer listener.Close()
	go func() {
		for {
			conn, err := listener.Accept()
			abort(err)
			if err == nil {
				s.connections.Set(conn.RemoteAddr().String(), conn)
			}
		}
	}()
}
func main() {

}
