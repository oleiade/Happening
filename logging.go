package happening

import (
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"os"
	"path/filepath"
)

// Log levels binding
var LogLevels = map[string]l4g.Level{
	l4g.DEBUG.String():    l4g.DEBUG,
	l4g.FINEST.String():   l4g.FINEST,
	l4g.FINE.String():     l4g.FINE,
	l4g.DEBUG.String():    l4g.DEBUG,
	l4g.TRACE.String():    l4g.TRACE,
	l4g.INFO.String():     l4g.INFO,
	l4g.WARNING.String():  l4g.WARNING,
	l4g.ERROR.String():    l4g.ERROR,
	l4g.CRITICAL.String(): l4g.CRITICAL,
}

// SetupLogger function ensures logging file exists, and
// is writable, and sets up a log4go filter accordingly
func SetupFileLogger(loggerName string, logLevel string, logFile string) error {
	dir := filepath.Dir(logFile)
	_, err := os.Stat(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Facteur can't open logging directory for writing. Reason: '%s'", err))
	}

	// check file permissions are correct
	_, err = os.Create(logFile)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERROR] Make sure %s is writable to the user launching facteur", dir))
	}

	l4g.AddFilter(loggerName,
		LogLevels[logLevel],
		l4g.NewFileLogWriter(logFile, false))

	return nil
}
