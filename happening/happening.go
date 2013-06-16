package main

import (
    // "fmt"
    "sync"
    happening "github.com/oleiade/happening"
)

func main() {
    var wg              sync.WaitGroup
    var store_channel   chan bool

    // Parse command line arguments
    cmdline := &happening.Cmdline{}
    cmdline.ParseArgs()

    // Launch client store listener routine
    wg.Add(1)
    client_store := happening.NewClientStore(store_channel, &wg)
    client_store.Run(*cmdline.Transport, *cmdline.ClientsPort)
}