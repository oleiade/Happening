package happening

import (
	"net"
)

func BuildTcpListener(transport string, host string, port string) (*net.TCPListener, error) {
	endpoint := host + port

	// Resolve tcp endpoint addr
	addr, err := net.ResolveTCPAddr(transport, endpoint)
	if nil != err {
		return nil, err
	}

	// Bind tcp socket, and set an unlimited timeout
	// FIXME: set timeout
	listener, err := net.ListenTCP(transport, addr)
	if err != nil {
		return nil, err
	}

	return listener, nil
}
