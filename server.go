package happening

import (
    l4g "github.com/alecthomas/log4go"
    "os"
    "os/signal"
    "syscall"
)

// Server implements the Service interface and exposes the Facteur
// different services.
type Server struct {
    Service
    EventsHandler        *EventsHandler
}

// Server initializes a new Server instance
func NewServer(handler *EventsHandler) *Server {
    return &Server{
        Service:        *NewService("Server"),
        EventsHandler:  handler,
    }
}

// Run launches the server's services and listens for SIGINT
// and SIGTERM signals to gracefully them on receive
func (s *Server) Run() error {
    defer s.waitGroup.Done()

    // Handle SIGINT and SIGTERM signals for gracefull shutdown sake.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        for sig := range ch {
            l4g.Logf(l4g.INFO, "[%s.Run] %s received, stopping the happening", s.name, sig)
            s.EventsHandler.NetworkService.Stop()
            os.Exit(1)
        }
    }()

    s.waitGroup.Add(2)
    s.EventsHandler.Serve()
    return nil
}

func ListenAndAcknowledge() error {
    return nil
}
