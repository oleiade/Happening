package happening

// Sockets buffer sizes
const (
    EVENTS_FLOW_BUF_SIZE = 4096
    CLIENTS_REG_BUF_SIZE = 256
)

// Messages constants
const (
    MSG_DELIMITER          = "\r\n"
    EVENT_PARAMS_SEPARATOR = '|'
)

// Timeouts in seconds
const (
    CLIENT_REG_CONN_TIMEOUT = 1
    CLIENT_REG_ID_TIMEOUT   = 30
    EVENT_REG_CONN_TIMEOUT  = 1
    EVENT_FLOW_TIMEOUT      = 30
)

// Configuration fallback constants
const (
    DEFAULT_CONFIG_FILE  = "/etc/happening/happening.conf"
    DEFAULT_LOG_FILE     = "/tmp/happening.log"
    DEFAULT_PID_FILE     = "/tmp/happening.pid"
    DEFAULT_TRANSPORT    = "tcp"
    DEFAULT_LOG_LEVEL    = "INFO"
    DEFAULT_DAEMON_MODE  = false
    DEFAULT_HOST         = "localhost"
    DEFAULT_CLIENTS_PORT = ":4044"
    DEFAULT_EVENTS_PORT  = ":4040"
)
