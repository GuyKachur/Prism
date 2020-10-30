package database

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/happierall/l"
	"github.com/pkg/errors"
)

type StorageAPI interface {
	SaveImage(model *Model) error
	LoadImage(filename string, parentid uint) (image.Image, error)
}

const ImageRoot = "./images/"

//SaveImage takes the image and saves it to the disk and returns at the path
func (instance *instance) SaveImage(model *Model) error {
	// if model.FileName = "" {
	// 	//file hasnt been saved yet, save in an intelligent way...
	// 	//if we do fileserver route, we can just use a directeroy
	// 	//i vote
	// 	// parent -> {filehash+ext}2523o52u5fsdk.jpg
	// 	// parent_folder filehash
	// 	//
	// }
	if ok := fileExists(model.FileName); ok {
		//file exists we already have it saved
		return nil
	}
	//i want to...
	//create a directory  of the model filename -minus extension
	//save the file to its filename
	//create /filename[-ext]/

	// //create the parent directory
	// if model.ParentID != 0 {
	// 	//if the parent isnt empty,
	// 	//
	// 	// createFile()
	// 	err := os.Mkdir(fmt.Sprintf("%d", model.ParentID), 0755)
	// 	if err != nil {
	// 		return errors.Wrap(err, "Unable to make parent directory")
	// 	}
	// }
	// file, err := os.Create(model.FileName)
	dirName := ImageRoot
	if model.ParentID != 0 {
		dirName = fmt.Sprintf("%s%d/", dirName, model.ParentID)
	}
	file, err := createFile(model.FileName, dirName)
	if err != nil {
		l.Error(errors.Wrap(err, "Error saving file: "+fmt.Sprintf("%d", model.UID)))
		return err
	}
	defer file.Close()
	img, _, err := image.Decode(bytes.NewReader(model.Image)) //forward
	if err != nil {
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s%d", ImageRoot, model.UID), os.ModePerm)
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

// func uidFromPath(fullPath string) string {
// 	_, file := path.Split(fullPath)
// 	uid := strings.Split(file, "-")[0]
// 	return uid
// }

func (instance *instance) LoadImage(filename string, parentID uint) (image.Image, error) {
	// dir, file := path.Split(filename)
	dirName := ImageRoot
	if parentID != 0 {
		dirName = fmt.Sprintf("%s%d/", dirName, parentID)
	}
	dirName = dirName + filename
	file, err := os.Open(filename)
	if err != nil {
		l.Error(errors.Wrap(err, "Error loading image at path: "+dirName))
		return nil, err
	}
	defer file.Close()
	//we should probably switch between the accepted outputs!
	image, _, err := image.Decode(file) //discard format name
	if err != nil {
		l.Error(errors.Wrap(err, "Error decoding image at path: "+dirName))
		return nil, err
	}
	return image, nil

}

func createFile(fileName, dir string) (*os.File, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		l.Debug("here")
		return nil, err
	}

	file, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return nil, err
	}
	return file, nil
}
