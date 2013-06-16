package happening

import "flag"

type Cmdline struct {
    DaemonMode      *bool
    Transport       *string
    ClientsPort     *string
    EventsPort      *string
}

func (c *Cmdline) ParseArgs() {
    c.DaemonMode = flag.Bool("d",
        DEFAULT_DAEMON_MODE,
        "Launches elevator as a daemon")
    c.Transport = flag.String("t",
        DEFAULT_TRANSPORT,
        "Sets the transport protocol to be used")
    c.ClientsPort = flag.String("c",
        DEFAULT_CLIENTS_PORT,
        "Port to be used for clients registration")
    c.EventsPort = flag.String("e",
        DEFAULT_EVENTS_PORT,
        "Port to be used for events registration")
    flag.Parse()
}