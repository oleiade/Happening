package happening

import "flag"

type Cmdline struct {
    ConfigFile      *string
    DaemonMode      *bool
    LogLevel        *string
    Transport       *string
    ClientsPort     *string
    EventsPort      *string
}

func (c *Cmdline) ParseArgs() {
    c.ConfigFile = flag.String("config",
        DEFAULT_CONFIG_FILE,
        "Specifies config file path")
    c.DaemonMode = flag.Bool("daemon",
        DEFAULT_DAEMON_MODE,
        "Launches elevator as a daemon")
    c.LogLevel = flag.String("log-level",
        DEFAULT_LOG_LEVEL,
        "Sets elevator verbosity")
    c.Transport = flag.String("transport",
        DEFAULT_TRANSPORT,
        "Sets the transport protocol to be used")
    c.ClientsPort = flag.String("clients-port",
        DEFAULT_CLIENTS_PORT,
        "Port to be used for clients registration")
    c.EventsPort = flag.String("events-port",
        DEFAULT_EVENTS_PORT,
        "Port to be used for events registration")
    flag.Parse()
}