package happening

import (
    "fmt"
    "net"
    "sync"
)

type ClientStore struct {
    Container   map[string]*net.IPConn
    Channel     chan bool
    group       *sync.WaitGroup
}

func NewClientStore(channel chan bool, group *sync.WaitGroup) *ClientStore {
    return &ClientStore {
        Container: make(map[string]*net.IPConn),
        Channel: channel,
        group: group,
    }
}


func (store *ClientStore) Run(netProto string, addr string) {
    ln, err := net.Listen(netProto, addr)
    if err != nil {
        // handle error
    }

    for {
        // If channel has been closed, set
        // sync as done, and goroutine ready to be
        // collected
        _, ok := <- store.Channel
        if !ok {
            store.group.Done()
            return
        }

        // Else, process the client connection
        conn, err := ln.Accept()
        if err != nil {
            continue
        }

        fmt.Println(conn)
    }
}
