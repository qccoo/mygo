package service

import (
    "context"
    "errors"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    bizservice "github.com/qccoo/w4/app/internal/biz/service"
    "github.com/gin-gonic/gin"
)

type Server struct {
    srv *http.Server
}

func NewServer(s bizservice.ImageService) Server {
    r := gin.Default()
    r.GET("/hello", func(c *gin.Context) {
        c.String(http.StatusOK, "hello!")
    })
    r.GET("/image", func(c *gin.Context) {
        c.String(http.StatusOK, s.GetImageAddr(c, "1"))
    })
    r.GET("/images/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.String(http.StatusOK, s.GetImageAddr(c, id))
    })
    return Server{srv: &http.Server{Addr: ":8080", Handler: r}}
}

func (srv *Server) Start() {
    go func() {
        if err := srv.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
            fmt.Printf("Listen with error: %s\n", err)
        }
    }()

    quit := make(chan os.Signal)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    fmt.Println("Shutting down server...")

    // To finish current requests
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.srv.Shutdown(ctx); err != nil {
        fmt.Printf("Server shutdown with error: %s\n", err)
    }
    fmt.Printf("Exiting...")
}

