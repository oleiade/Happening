package happening

import (
    "bytes"
    "errors"
    "fmt"
    l4g "github.com/alecthomas/log4go"
    "net"
    "strconv"
    "strings"
    "time"
)

// Event represents a event sent by the source
type Event struct {
    Timestamp    int
    Type         string
    Emitter      string
    Weight       string
    raw          string
}

// EventsHandler is a NetworkService implementation which
// exposes a PriorityQueue to store incoming events,
// a Buffer to bufferize events flow socket input,
// and can notify other services listening on it's EventsChannel
// of events received and ready to process.
type EventsHandler struct {
    NetworkService
    Queue         *Queue
    Buffer        []byte
    EventsChannel chan *Event
}

// NewEventsHandler initializes an EventsHandler.
// Note that it calls the heap interface Init method
// over PriorityQueue so you won't have to do it later.
func NewEventsHandler() *EventsHandler {
    return &EventsHandler{
        NetworkService: *NewNetworkService("EventsHandler"),
        Queue:          &Queue{},
        EventsChannel:  make(chan *Event),
    }
}

// NewEvent initializes an event from it's component
func NewEvent(ts int, typ string, emitter string, weight string) *Event {
    return &Event{
        Timestamp: ts,
        Type: typ,
        Emitter: emitter,
        Weight: weight,
    }
}

// NewEventFromRaw initializes an Event from it's raw representation
// which MSG_DELIMITER has been removed from.
func NewEventFromRaw(raw string) (*Event, error) {
    event := new(Event)
    err := event.FromRaw(raw)
    if err != nil {
        return nil, err
    }

    return event, nil
}

// FromRaw instantiates and returns an initialized Event
// from it's raw description. As event flow splits event messages
// based on the MSG_DELIMITER and discards it, the method artificially
// restores the MSG_DELIMITER in Event.raw attribute.
func (e *Event) FromRaw(raw string) error {
    e.raw = raw + MSG_DELIMITER // Keep track of the raw version with MSG_DELIMITER
    parts := strings.Split(strings.Trim(raw, MSG_DELIMITER), string(EVENT_PARAMS_SEPARATOR))

    if len(parts) >= 2 {
        // Set up event sequence id or fail
        ts, err := strconv.Atoi(parts[0])
        if err != nil {
            return errors.New(fmt.Sprintf("[%s.FromRaw] Can't convert %s to int", "Event", parts[0]))
        }

        e.Timestamp = ts
        e.Type = parts[1]

        if len(parts) >= 3 {
            e.Emitter = parts[2]
        }

        if len(parts) >= 4 {
            e.Weight = parts[3]
        }
    }

    return nil
}

// Serve should be run as a long-running goroutine.
// It runs a HandleConnexion long-running goroutine, awaits for
// new connexions sent through new_sources channel, and eventually
// starts an HandleEvents goroutine to process new incoming events.
func (m *EventsHandler) Serve() {
    defer m.waitGroup.Done()
    l4g.Info("Events source ready for the flow")

    // Channels dedicated to handle the sources and events
    // stream processing goroutines state
    source_state := make(chan bool)
    events_state := make(chan bool)

    new_sources := make(chan *net.TCPConn)
    m.waitGroup.Add(1)
    go m.HandleConnexion(source_state, new_sources)

    for {
        select {
        // If channel has been closed, or a shutdown
        // signal has been sent, set sync as done
        // and goroutine ready to be collected
        case <-m.ch:
            close(source_state)
            close(events_state)
            return
        // Otherwise, process the events source connection and events
        case new_source := <-new_sources:
            m.waitGroup.Add(1)
            go m.HandleEvents(events_state, new_source)
        }
    }
}

// HandleConnexion should be used as a long-running goroutine to listen
// on EventsHandler socket for new event source connexions.
// Each new connexion will be sent back to HandleConnexion caller through
// sources channel.
// Anytime HandleConnexion can be stoppped by closing it's state channel.
func (m *EventsHandler) HandleConnexion(state chan bool, sources chan *net.TCPConn) {
    defer m.waitGroup.Done()

    for {
        select {
        case <-state:
            return
        default:
            // Awainting for the events source to connect
            m.Socket.SetDeadline(time.Now().Add(time.Duration(EVENT_REG_CONN_TIMEOUT) * time.Second))
            source, err := m.Socket.AcceptTCP()
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                    continue
                }
                l4g.Error(fmt.Sprintf("[%s.HandleConnexion] %s", m.name, err))
                return
            }

            sources <- source
        }
    }
}

// HandleEvents should be run as a long-running goroutine to listen
// on the EventsHandler source and process incoming events.
// It reads on the socket, maintain the EventsHandler buffer,
// extracts the message from the buffer, instantiates Events,
// and pushes them to the PriorityQueue.
func (m *EventsHandler) HandleEvents(events_state chan bool, source *net.TCPConn) {
    defer m.waitGroup.Done()
    defer source.Close()

    for {
        select {
        case <-events_state:
            return
        default:
            socket_input := make([]byte, EVENTS_FLOW_BUF_SIZE)

            source.SetDeadline(time.Now().Add(time.Duration(EVENT_FLOW_TIMEOUT) * time.Second))
            // Await for client id to be sent
            read_len, err := source.Read(socket_input)
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                    continue
                }
                l4g.Error(fmt.Sprintf("[%s.HandleEvents] Events source connexion closed", m.name))
                return
            }

            // In order to protect the events splitted accross two
            // socket read buffers, we copy the eventual rest and the
            // newly received data into a new buffer removing zero
            // bytes from the actual socket buffer
            data := make([]byte, read_len+len(m.Buffer))
            copy(data, m.Buffer)
            copy(data[len(m.Buffer):], socket_input[:read_len])
            m.Buffer = []byte{} // reinitialize to empty

            // Was the socket buffer ended with an incomplete
            // event message? And was the actual content of buffer
            // ended with the msg delimiter?
            incomplete := incompleteEventMessage(socket_input)
            backslash_ended := backslashEndedEventMessage(data)

            // extract events from the input data
            items := strings.Split(string(data), MSG_DELIMITER)

            // If socket buffer was ended with an incomplete message
            // push the rest in the EventsHandler buffer, and remove
            // last event extracted as it is incomplete.
            // Else, if the last event found was actually ended with
            // MSG_DELIMITER, strings.Split will return an empty string
            // elem after it, so let's remove it.
            if incomplete {
                m.Buffer = []byte(items[len(items)-1])
                items = items[:len(items)-1]
            } else if backslash_ended {
                items = items[:len(items)-1]
            }

            m.PushEventsToQueue(items)
        }
    }
}

// PushEventsToQueue adds a list of Event instances to
// the EventsHandler PriorityQueue.
func (m *EventsHandler) PushEventsToQueue(events []string) {
    for _, item := range events {
        event, err := NewEventFromRaw(item)
        if err != nil {
            l4g.Error(fmt.Sprintf("[%s.PushEventsToQueue] %s", m.name, err))
            continue
        }

        l4g.Info(fmt.Sprintf("[%s.PushEventsToQueue] %.2d:%q inserted in queue", m.name, event.Weight, event.raw))
        m.Queue.Push(event)

        // Send event in events channel for listener
        // to be notified of the event pushed to queue
        // the forwarder for example.
        go func() { m.EventsChannel <- event }()
    }
}

// incompleteEventMessage checks if the EventsHandler input socket
// read message is incomplete (non-terminated by MSG_DELIMITER)
func incompleteEventMessage(socket_buffer []byte) bool {
    return !bytes.HasSuffix(socket_buffer, []byte(MSG_DELIMITER)) && socket_buffer[len(socket_buffer)-1] != byte(0)
}

func backslashEndedEventMessage(data []byte) bool {
    return bytes.Compare(data[len(data)-len(MSG_DELIMITER):], []byte(MSG_DELIMITER)) == 0
}
