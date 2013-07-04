package happening

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"time"
	l4g "github.com/alecthomas/log4go"
)

// ClientStore implements the NetworkService interface and exposes
// a Container to store client id to client socket relations (in order
// to be able to communicate with registered users).
// Nota:
//     * As golang maps are not safe for concurrent use,
// ClientStore implements the sync.RWMutex interface
type ClientStore struct {
	sync.RWMutex
	NetworkService
	Container     map[string]*net.TCPConn
}

// NewClientStore initializes a new ClientStore
func NewClientStore() *ClientStore {
	return &ClientStore{
		NetworkService: *NewNetworkService("ClientStore"),
		Container:      make(map[string]*net.TCPConn),
	}
}

// Serve runs a HandleConnexions long-running goroutine, awaits for
// new connexions sent through new_connexions channel, and eventually
// starts an HandleClient goroutine to go trough new client connexion
// registration process
func (store *ClientStore) Serve() {
	defer store.waitGroup.Done()
	l4g.Info("Clients store registration routine started")

	// Start the clients connexions handling goroutine and bind
	// an input channel to it in order to be able to ask it
	// to stop gracefully, as well as an output channel to
	// retrieve the connexions it retrieves
	connexions_state := make(chan bool)
	new_connexions := make(chan *net.TCPConn)
	store.waitGroup.Add(1)
	go store.HandleConnexions(connexions_state, new_connexions)

	for {
		select {
		// If channel has been closed, or a shutdown
		// signal has been sent, set sync as done
		// and goroutine ready to be collected
		case <-store.ch:
			close(connexions_state)
			return
		// If a new client connexion has been sent by
		// the dedicated goroutine, launch the client
		//registration goroutine
		case new_connexion := <-new_connexions:
			l4g.Logf(l4g.DEBUG, "[%s.Server] Client connection received from: %s", store.name, new_connexion.RemoteAddr())
			store.waitGroup.Add(1)
			go store.HandleClient(new_connexion)
		}
	}
}

// HandleConnexions should be used as a long-running goroutine to listen
// on ClientStore's socket for new connexions.
// Each new connexion will be sent back to HandleConnexions caller through
// new_connexions channel.
// Anytime HandleConnexions can be stoppped by sending a boolean value
// through it's status channel
func (store *ClientStore) HandleConnexions(status chan bool, new_connexions chan *net.TCPConn) {
	defer store.waitGroup.Done()

	for {
		select {
		case <-status:
			return
		default:
			// Await on incoming connection and spawn a new goroutine
			store.Socket.SetDeadline(time.Now().Add(time.Duration(CLIENT_REG_CONN_TIMEOUT) * time.Second))
			conn, err := store.Socket.AcceptTCP()
			if err != nil {
				// If AcceptTCP timeouted, jump to next iteration, in order
				// to re-check the status channel, and eventually re-try
				// the AcceptTCP
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				l4g.Logf(l4g.ERROR, "[%s.HandleConnexion] %s", store.name, err)
			}
			new_connexions <- conn
		}
	}
}

// HandleClient listens on a newly connected client socket for
// incoming client's id, and eventually registers it into ClientStore's
// sockets container.
func (store *ClientStore) HandleClient(conn *net.TCPConn) {
	defer store.waitGroup.Done()

	buf := make([]byte, CLIENTS_REG_BUF_SIZE)

	// Set client id receiving timeout
	conn.SetDeadline(time.Now().Add(time.Duration(CLIENT_REG_ID_TIMEOUT) * time.Second))
	// Await for client id to be sent
	if _, err := conn.Read(buf); err != nil {
		l4g.Error(err)
		return
	}

	// Parse client id from raw message
	scanner := bufio.NewScanner(strings.NewReader(string(buf)))
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	id := scanner.Text()

	// Store client connexion
	store.Lock()
	store.Container[id] = conn
	store.Unlock()

	l4g.Logf(l4g.DEBUG, "[%s.HandleClient] Client with id %s registered", store.name, id)
}
