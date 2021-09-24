package main

import (
    "fmt"
    "os"
)

func main() {
    srv, err := InitializeServer()
    if err != nil {
        fmt.Printf("Failed to InitializeServer: %s\n", err)
        os.Exit(2)
    }
    srv.Start()
}
