package main

import (
    "os"
    "fmt"
    "log"
    "syscall"
    "os/signal"
    l4g "github.com/alecthomas/log4go"
    happening "github.com/oleiade/happening"
)

func main() {
    var err             error

    l4g.Info("Happening started")

    // Parse command line arguments
    cmdline := &happening.Cmdline{}
    cmdline.ParseArgs()

    // Load configuration from file
    config := happening.NewConfig()
    err = config.FromFile(*cmdline.ConfigFile, "core")
    if err != nil {
        log.Fatal(err)
    }
    config.UpdateFromCmdline(cmdline)

    // Set up loggers
    l4g.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
    err = happening.SetupFileLogger("file", config.LogLevel, config.LogFile)
    if err != nil {
        log.Fatal(err)
    }

    // build client store
    events_handler := happening.NewEventsHandler()
    events_handler.InitSocket(*cmdline.Host, *cmdline.EventsPort)

    // build server
    server := happening.NewServer(events_handler)

    l4g.Info("Hapening events listener routine started")

    // Handle SIGINT and SIGTERM.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        for sig := range ch {
            l4g.Info(fmt.Sprintf("%s received, stopping the happening", sig))
            server.Stop()
            os.Exit(1)
        }
    }()

    server.Run()
}