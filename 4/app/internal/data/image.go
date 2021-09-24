package data

import (
	"context"
	// "database/sql"
	"fmt"
	"regexp"
	"github.com/qccoo/w4/app/internal/biz/model"
)

type ImageRepoImpl struct {
	// db *sql.DB
}

func NewRepo() ImageRepoImpl {
	return ImageRepoImpl{}
}

// TODO: Add and handle DB

var (
	reg, _ = regexp.Compile("[^0-9]+")
)

func (im *ImageRepoImpl) GetImage(ctx context.Context, id string) (*model.Image, error) {
	// img := &model.Image{}
	// row := d.db.QueryRow("select id, text, file from images where id = ?", id)
	// err := row.Scan(&img.ID, &img.Text, &img.File)
	// if err != nil {
	// 	return nil, err
	// }
	// return img, nil
    proId := reg.ReplaceAllString(id, "")
    trunc := "1"
    if (len(proId)) > 0 {
    	trunc = proId[:1]
    }
	file := fmt.Sprintf("free_img%v.jpeg", trunc)
	return &model.Image{Id: "id" + trunc, Description: "fake image " + trunc, File: file}, nil
}


func (im *ImageRepoImpl) ListImages(ctx context.Context) ([]*model.Image, error) {
	var images []*model.Image
	images = append(images, &model.Image{Id: "id1", Description: "fake image 1", File: "free_img1.jpeg"})
	images = append(images, &model.Image{Id: "id2", Description: "fake image 2", File: "free_img2.jpeg"})
	return images, nil
}

