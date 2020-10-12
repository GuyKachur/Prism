package database

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"

	"github.com/happierall/l"
	"github.com/pkg/errors"
)

type StorageAPI interface {
	SaveImage(model *Model) error
	LoadImage(path string) (image.Image, string, error)
}

const ImageRoot = "/images/"

//SaveImage takes the image and saves it to the disk and returns at the path
func (instance *instance) SaveImage(model *Model) error {
	if ok := fileExists(model.FileName); ok {
		//file exists we already have it saved
		return nil
	}
	file, err := os.Create(model.FileName)
	if err != nil {
		l.Error(errors.Wrap(err, "Error saving file: "+fmt.Sprintf("%d", model.UID)))
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(bytes.NewReader(model.Image)) //forward
	if err != nil {
		return err
	}
	return jpeg.Encode(file, img, &jpeg.Options{95}) //default

}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func uidFromPath(fullPath string) string {
	_, file := path.Split(fullPath)
	uid := strings.Split(file, "-")[0]
	return uid
}

func (instance *instance) LoadImage(path string) (image.Image, string, error) {
	uid := uidFromPath(path)
	file, err := os.Open(path)
	if err != nil {
		l.Error(errors.Wrap(err, "Error loading image at path: "+path))
		return nil, uid, err
	}
	defer file.Close()
	image, _, err := image.Decode(file) //discard format name
	if err != nil {
		l.Error(errors.Wrap(err, "Error loading image at path: "+path))
		return nil, uid, err
	}
	return image, uid, nil

}
