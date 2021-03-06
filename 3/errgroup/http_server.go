package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "golang.org/x/sync/errgroup"
)

func main() {
    g, ctx := errgroup.WithContext(context.Background())
    server := &http.Server{Addr: ":8080"}
    http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(resp, "Hello!")
    })

    g.Go(func() error {
        fmt.Println("Server(:8080)...")
        return server.ListenAndServe()
    })

    g.Go(func() error {
        <-ctx.Done()
        fmt.Println("Http server context done. Shutting down...")
        return server.Shutdown(context.TODO())
    })

    // Handle signals
    g.Go(func() error {
        signals := make(chan os.Signal, 1)
        signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Signals handler context done.")
                return ctx.Err()
            case s := <-signals:
                // Returns a non-nil error to cancel the derived Context
                return fmt.Errorf("Received exit signal: %v", s)
            }
        }
    })

    err := g.Wait()
    if err != nil {
        fmt.Println("Exited with error:", err)
    } else {
        fmt.Println("All done without errors.")
    }
}
