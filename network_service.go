package happening

import (
    "net"
    "time"
    "fmt"
    l4g "github.com/alecthomas/log4go"
)

// NetworkService is built on the Service structure and adds the support
// for a net.TCPListener in order to create a service ready for networking.
type NetworkService struct {
    Service
    Socket              *net.TCPListener
    ConnexionsLifeline  chan bool
    IncomingConnexions  chan *net.TCPConn
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
        ConnexionsLifeline: make(chan bool),
        IncomingConnexions: make(chan *net.TCPConn),
    }
    return ns
}

// InitSocket is a helper method to set up a NetworkService instance
// socket from a transport, host and port strings.
func (ns *NetworkService) initSocket(host string, port string) error {
    socket, err := BuildTcpListener("tcp", host, port)
    if err != nil {
        return err
    }

    ns.Socket = socket

    return nil
}

func (ns *NetworkService) Start(host string, port string) (err error) {
    err = ns.initSocket(host, port)
    if err != nil {
        return err
    }

    ns.waitGroup.Add(1)
    go ns.HandleConnexions()

    return nil
}

// Stop the NetworkService by closing the service's channel and socket.
// Blocks until the network service is really stopped.
func (ns *NetworkService) Stop() {
    close(ns.ch)
    ns.waitGroup.Wait()
    ns.Socket.Close()
}

// HandleConnexion should be used as a long-running goroutine to listen
// on NetworkService socket for new event source connexions.
// Each new connexion will be sent back to HandleConnexion caller through
// sources channel.
// Anytime HandleConnexion can be stoppped by closing it's lifeline channel.
func (ns *NetworkService) HandleConnexions() {
    defer ns.waitGroup.Done()

    for {
        select {
        case <- ns.ConnexionsLifeline:
            return
        default:
            // Awainting for the events source to connect
            ns.Socket.SetDeadline(time.Now().Add(time.Duration(EVENT_REG_CONN_TIMEOUT) * time.Second))
            source, err := ns.Socket.AcceptTCP()
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                    continue
                }
                l4g.Error(fmt.Sprintf("[%s.HandleConnexion] %s", ns.name, err))
                return
            }

            ns.IncomingConnexions <- source
        }
    }
}

