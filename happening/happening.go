package main

import (
    "fmt"
    happening "github.com/oleiade/happening"
)

func main() {
    // Parse command line arguments
    cmdline := &happening.Cmdline{}
    cmdline.ParseArgs()

    fmt.Println(cmdline)
}