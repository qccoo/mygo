package main

import (
    "github.com/google/wire"

    bizservice "github.com/qccoo/w4/app/internal/biz/service"
    "github.com/qccoo/w4/app/internal/data"
    "github.com/qccoo/w4/app/internal/service"
)

func InitializeServer() (service.Server, error) {
    wire.Build(service.NewServer, bizservice.NewImageService, data.NewRepo)
    return service.Server{}, nil
}
