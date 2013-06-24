package main

import (
    "os"
    "fmt"
    "log"
    "time"
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

    // Launch client store listener routine
    l4g.Info("Happening clients registration routine started")
    client_store := happening.NewClientStore()
    go client_store.Serve(*cmdline.Transport, *cmdline.ClientsPort)

    l4g.Info("Hapening events listener routine started")

    // Handle SIGINT and SIGTERM.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        for sig := range ch {
            l4g.Info(fmt.Sprintf("%s received, stopping the happening", sig))
            client_store.Stop()
            l4g.Info("Client store stopped")
            os.Exit(1)
        }
    }()

    for {
        time.Sleep(100)
    }
}