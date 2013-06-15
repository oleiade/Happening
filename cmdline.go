package happening

import "flag"

type Cmdline struct {
    DaemonMode *bool
}

func (c *Cmdline) ParseArgs() {
    c.DaemonMode = flag.Bool("d",
        DEFAULT_DAEMON_MODE,
        "Launches elevator as a daemon")
    flag.Parse()
}