package happening

import (
    "bytes"
    "fmt"
    l4g "github.com/alecthomas/log4go"
    "net"
    "strings"
    "time"
)

type EventsHandler struct {
    NetworkService
    Queue         *Queue
    Buffer        []byte
    EventsChannel chan *Event
}

// NewEventsHandler initializes an EventsHandler.
func NewEventsHandler() *EventsHandler {
    return &EventsHandler{
        NetworkService: *NewNetworkService("EventsHandler"),
        Queue:          NewQueue(EVENTS_QUEUE_SIZE),
        EventsChannel:  make(chan *Event),
    }
}

// Serve should be run as a long-running goroutine.
// It runs a HandleConnexion long-running goroutine, awaits for
// new connexions sent through new_sources channel, and eventually
// starts an HandleEvents goroutine to process new incoming events.
func (m *EventsHandler) Serve() {
    defer m.waitGroup.Done()
    l4g.Info("Events source ready for the flow")
    events_state := make(chan bool)

    for {
        select {
        // If channel has been closed, or a shutdown
        // signal has been sent, set sync as done
        // and goroutine ready to be collected
        case <- m.ch:
            close(m.ConnexionsLifeline)
            close(events_state)
            return
        // Otherwise, process the events source connection and events
        case new_source := <- m.IncomingConnexions:
            m.waitGroup.Add(1)
            go m.HandleEvents(events_state, new_source)
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

            items := m.ExtractEventsFromSocketInput(socket_input, read_len)
            go m.PushEventsToQueue(items)
        }
    }
}

func (m *EventsHandler) ExtractEventsFromSocketInput(input []byte, read_len int) []string {
    // In order to protect the events splitted accross two
    // socket read buffers, we copy the eventual rest and the
    // newly received data into a new buffer removing zero
    // bytes from the actual socket buffer
    data := make([]byte, read_len+len(m.Buffer))
    copy(data, m.Buffer)
    copy(data[len(m.Buffer):], input[:read_len])
    m.Buffer = []byte{} // reinitialize to empty

    // Was the socket buffer ended with an incomplete
    // event message? And was the actual content of buffer
    // ended with the msg delimiter?
    is_incomplete := !bytes.HasSuffix(input, []byte(MSG_DELIMITER)) && input[len(input)-1] != byte(0)
    is_backslash_ended := bytes.Compare(data[len(data)-len(MSG_DELIMITER):], []byte(MSG_DELIMITER)) == 0

    // extract events from the input data
    items := strings.Split(string(data), MSG_DELIMITER)

    // If socket buffer was ended with an incomplete message
    // push the rest in the EventsHandler buffer, and remove
    // last event extracted as it is incomplete.
    // Else, if the last event found was actually ended with
    // MSG_DELIMITER, strings.Split will return an empty string
    // elem after it, so let's remove it.
    if is_incomplete {
        m.Buffer = []byte(items[len(items)-1])
        items = items[:len(items)-1]
    } else if is_backslash_ended {
        items = items[:len(items)-1]
    }

    return items
}

// PushEventsToQueue adds a list of Event instances to
// the EventsHandler PriorityQueue.
func (m *EventsHandler) PushEventsToQueue(events []string) {
    for _, event := range events {
        event, err := NewEventFromRaw(event)
        if err != nil {
            l4g.Error(fmt.Sprintf("[%s.PushEventsToQueue] %s", m.name, err))
            continue
        }

        l4g.Info(fmt.Sprintf("[%s.PushEventsToQueue] %s inserted in queue", m.name, event))
        m.Queue.Push(event)

        // Send event in events channel for listener
        // to be notified of the event pushed to queue
        // the forwarder for example.
        go func() { m.EventsChannel <- event }()
    }
}
