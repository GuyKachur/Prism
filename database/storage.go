package database

import (
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
	SaveImage(image image.Image, model Model) (string, error)
	LoadImage(path string) (image.Image, string, error)
}

const imageRoot = "/images/"

//SaveImage takes the image and saves it to the disk at the path
func (instance *instance) SaveImage(image image.Image, model Model) (string, error) {
	path := imageRoot + fmt.Sprintf("%d", model.UID) + "-" + model.FileName
	// if ok := fileExists(path); ok {
	file, err := os.Create(path)
	if err != nil {
		l.Error(errors.Wrap(err, "Error saving file: "+fmt.Sprintf("%d", model.UID)))
		return "", err
	}
	defer file.Close()
	return path, jpeg.Encode(file, image, &jpeg.Options{95}) //default

	// } else {
	// 	return
	// }

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
