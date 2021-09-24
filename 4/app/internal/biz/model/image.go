package model


import (
	"context"
)

type Image struct {
	Id   string `json:"id"`
	Description string `json:"description"`
	File string `json:"file"`
}

type ImageRepo interface {
	GetImage(ctx context.Context, id string) (*Image, error)
	ListImages(ctx context.Context) ([]*Image, error)
	// AddImage(ctx context.Context, image *Image) error
}
