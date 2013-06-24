package happening

import (
	"fmt"
	"net"
	"time"
	l4g "github.com/alecthomas/log4go"
)

type ClientStore struct {
	Service
	Container map[string]*net.IPConn

}

func NewClientStore() *ClientStore {
	return &ClientStore{
		Service: *NewService(),
		Container: make(map[string]*net.IPConn),
	}
}

func (store *ClientStore) Serve(netProto string, endpoint string) {
	defer store.waitGroup.Done()

	// Resolve tcp endpoint addr
	addr, err := net.ResolveTCPAddr(netProto, endpoint)
	if nil != err {
		l4g.Critical(err)
	}

	// Bind tcp socket, and set an unlimited timeout
	// FIXME: set timeout
	ln, err := net.ListenTCP(netProto, addr)
	if err != nil {
		l4g.Critical(err)
	}

	for {
		select {
			// If channel has been closed, or a shutdown
			// signal has been sent, set sync as done
			// and goroutine ready to be collected
			case <-store.ch:
				ln.Close()
				l4g.Info("Stopping clients store")
				return
			// Otherwise, process the client connection
			default:
		}

		ln.SetDeadline(time.Now().Add(1e9))
		conn, err := ln.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			l4g.Error(err)
		}

		l4g.Info(fmt.Sprintf("Connection received from: %s", conn.RemoteAddr()))
		_, err = conn.Write([]byte("Ground control reiceived your message major Tom"))
		if err != nil {
			return
		}
	}
}
