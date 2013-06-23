package main

import (
    "log"
    "sync"
    l4g "github.com/alecthomas/log4go"
    happening "github.com/oleiade/happening"
)

func main() {
    var err             error
    var wg              sync.WaitGroup
    var store_channel   chan bool

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

    l4g.Info("Happening clients registration routine started")
    // Launch client store listener routine
    wg.Add(1)
    client_store := happening.NewClientStore(store_channel, &wg)
    client_store.Run(*cmdline.Transport, *cmdline.ClientsPort)

    l4g.Info("Hapening events listener routine started")
}