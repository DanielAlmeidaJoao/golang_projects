package commons

import "net"

type DataWrapper struct {
	From net.Addr
	Data []byte
}
