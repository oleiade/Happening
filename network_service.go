package happening

import (
    "net"
)

// NetworkService is built on the Service structure and adds the support
// for a net.TCPListener in order to create a service ready for networking.
type NetworkService struct {
    Service
    Socket *net.TCPListener
}

// NewNetworkService builds a new NetworkService instance. In order
// to avoid possible errors at initialization (and not to pollute the
// initializer) Service attribute is intentionaly set to nil.
// Please use the InitSocket method over a NetworkService to setup it's
// socket.
func NewNetworkService(name string) *NetworkService {
    ns := &NetworkService{
        Service: *NewService(name),
        Socket:  nil,
    }
    return ns
}

// InitSocket is a helper method to set up a NetworkService instance
// socket from a transport, host and port strings.
func (ns *NetworkService) InitSocket(host string, port string) error {
    socket, err := BuildTcpListener("tcp", host, port)
    if err != nil {
        return err
    }

    ns.Socket = socket

    return nil
}

// Stop the NetworkService by closing the service's channel and socket.
// Blocks until the network service is really stopped.
func (ns *NetworkService) Stop() {
    close(ns.ch)
    ns.waitGroup.Wait()
    ns.Socket.Close()
}
