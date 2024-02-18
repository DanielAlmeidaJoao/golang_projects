package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	defer ln.Close()
	fmt.Println("SERVER LISTENING...")
	if err != nil {
		return err
	}
	s.ln = ln
	go s.acceptLoop()
	fmt.Println("SERVER WAITING FOR SHUTDOWN...")
	<-s.quitch
	fmt.Println("SERVER shutting DOWN...")
	time.Sleep(1)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("ACCEPT ERROR:", err)
			return
		}
		fmt.Printf("SERVER ACCEPTED A CONNECTION. %s\n", conn.RemoteAddr().String())
		fmt.Println("SERVER MONITORING THE CONNECTION FOR INCOMING DATA...")
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("READ ERROR:", err)
			continue
		}
		msg := buf[:n]
		strMsg := strings.Replace(string(msg), "\r\n", "", -1)
		if strings.Compare(strMsg, "bye") == 0 {
			s.quitch <- struct{}{}
		}

		fmt.Printf("< %s > ]] \n", strMsg)
	}
}
func main() {

	server := NewServer(":3000")
	server.start()

}
