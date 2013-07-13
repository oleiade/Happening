package happening

import "flag"

type Cmdline struct {
	DaemonMode *bool
	ConfigFile *string
	PidFile    *string
	LogFile    *string
	LogLevel   *string
	Host       *string
	EventsPort *string
}

func (c *Cmdline) ParseArgs() {
	c.DaemonMode = flag.Bool("daemon",
		DEFAULT_DAEMON_MODE,
		"Launches elevator as a daemon")
	c.ConfigFile = flag.String("config",
		DEFAULT_CONFIG_FILE,
		"Specifies config file path")
	c.PidFile = flag.String("pid-file",
		DEFAULT_PID_FILE,
		"Specifies which pid file happening should maintain when in daemon mode")
	c.LogFile = flag.String("log-file",
		DEFAULT_LOG_FILE,
		"Specifies in which file happening should eventually log")
	c.LogLevel = flag.String("log-level",
		DEFAULT_LOG_LEVEL,
		"Sets elevator verbosity")
	c.Host = flag.String("host",
		DEFAULT_HOST,
		"Sets the host to bind happening sockets to")
	c.EventsPort = flag.String("events-port",
		DEFAULT_EVENTS_PORT,
		"Port to be used for events registration")
	flag.Parse()
}
