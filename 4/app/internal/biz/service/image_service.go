package service


import (
	"fmt"

	"github.com/qccoo/w4/app/internal/data"
	"github.com/qccoo/w4/app/internal/pkg/converters"
    "github.com/gin-gonic/gin"
)

type ImageService struct {
	repo data.ImageRepoImpl
}

func NewImageService(repo data.ImageRepoImpl) ImageService {
	return ImageService{repo: repo}
}

func (s *ImageService) GetImageAddr(c *gin.Context, id string) string {
	img, err := s.repo.GetImage(c, id)
	if err != nil {
        fmt.Printf("Failed to GetImageAddr: %s\n", err)
		return ""
	}
	return converters.ConvertImageAddr(img.File)
}