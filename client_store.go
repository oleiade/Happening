package happening

import (
	"net"
	"sync"
)

type ClientStore struct {
	Container map[string]*net.IPConn
	Channel   chan bool
	group     *sync.WaitGroup
}

func NewClientStore(channel chan bool, group *sync.WaitGroup) *ClientStore {
	return &ClientStore{
		Container: make(map[string]*net.IPConn),
		Channel:   channel,
		group:     group,
	}
}

func (store *ClientStore) Run(netProto string, addr string) {
	ln, err := net.Listen(netProto, addr)
	if err != nil {
		// handle error
	}

	for {
		select {
		// If channel has been closed, or a shutdown
		// signal has been sent, set sync as done
		// and goroutine ready to be collected
		case _ = <-store.Channel:
			store.group.Done()
			return
		// Otherwise, process the client connection
		default:
			conn, err := ln.Accept()
			if err != nil {
				continue
			}

			_, err = conn.Write([]byte("Ground control reiceived your message major Tom"))
			if err != nil {
				return
			}
		}
	}
}
