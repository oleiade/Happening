package happening

// Sockets buffer sizes
const (
	EVENTS_FLOW_BUF_SIZE = 4096
)

// Storage backends constants
const (
	LEVELDB_LRU_CACHE_SIZE = 64 * 1048576 // 64Mo
)

// Messages constants
const (
	MSG_DELIMITER          = "\r\n"
	EVENT_PARAMS_SEPARATOR = '|'
)

// Timeouts in seconds
const (
	EVENT_REG_CONN_TIMEOUT = 1
	EVENT_FLOW_TIMEOUT     = 30
)

// Internal events queue size
const (
	EVENTS_QUEUE_SIZE = 4096
)

// Configuration fallback constants
const (
	DEFAULT_CONFIG_FILE  = "/etc/happening/happening.conf"
	DEFAULT_STORAGE_PATH = "/tmp"
	DEFAULT_LOG_FILE     = "/tmp/happening.log"
	DEFAULT_PID_FILE     = "/tmp/happening.pid"
	DEFAULT_TRANSPORT    = "tcp"
	DEFAULT_LOG_LEVEL    = "INFO"
	DEFAULT_DAEMON_MODE  = false
	DEFAULT_HOST         = "localhost"
	DEFAULT_EVENTS_PORT  = ":4040"
)
