package server

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"refract/database"
	"refract/refract"

	"github.com/go-chi/chi"
	"github.com/happierall/l"
)

//assumes Image is alread in database
//This will need to be split, if url is full -> go url route, if not go database route

func RefractHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid != "" {
		HandleError(w, errors.New("Missing Parameter UID"))
		return
	}

	//get image and metadata
	model, err := database.Instance.GetImage(uid)
	if err != nil {
		HandleError(w, err)
		return
	}

	ca := r.Header.Get("config")
	config := &refract.Config{} //everything minus the input and the output
	err = json.Unmarshal([]byte(ca), config)
	if err != nil {
		HandleError(w, err)
		return
	}
	// save image and spit path out into config
	//save image to filesystem for primitives use
	err = database.Instance.SaveImage(model)
	if err != nil {
		HandleError(w, err)
		return
	}
	config.Input = model.FileName
	//where do we want to solve the resulting image?
	// i want to make it parentuid-name-file
	// so check incoming filename for path in os...
	newPath := database.ImageRoot + "temp-" + model.FileName
	l.Debug("New image being created at ", newPath)
	config.Output = newPath
	// config.Output = database.ImageRoot + config.Name + fmt.Sprintf("parent-%d", model.UID) + ""

	// //Image will be duplicated to storage...
	// config.output = newPath
	// //start building new image

	//heres what i know, i have a model and an image. And primitive is looking to take that image and write it tothe output
	commandByteArray, _, err := refract.Primitive(*config)
	l.Debug(commandByteArray)
	if err != nil {
		HandleError(w, err)
		return
	}
	name := r.Header.Get("name")
	if name != "" {
		name = config.Name + "-" + model.Name
	}
	// imgFile, err := os.Open(newPath)
	// defer imgFile.Close()
	// if err != nil {
	// 	HandleError(w, err)
	// 	return
	// }

	// img, _, err := image.Decode(imgFile)
	// if err != nil {
	// 	HandleError(w, err)
	// 	return
	// }

	tags := r.Header.Get("tags")
	img, err := ioutil.ReadFile(newPath)
	if err != nil {
		HandleError(w, err)
		return
	}
	newHash := md5.Sum(img)
	newModel := &database.Model{
		Name:     name,
		Image:    img,
		FileName: config.Name + "-" + fmt.Sprintf("parent-%d", model.UID),
		ParentID: model.UID,
		Tags:     tags,
		FileHash: newHash[:],
	}

	//save resulting image
	err = database.Instance.Upload(newModel)
	if err != nil {
		HandleError(w, err)
		return
	}
	//that will hopefully return a success.... though it might only return on error state
	//fill the output path
	//now we need to grab the outputted image,
	// newImage := &Model{
	// 	Name: name
	// }

	//model
	// 	UID       uint      `gorm:"primaryKey" json:"uid,omitempty"`
	// 	CreatedAt time.Time `json:"created_at,omitempty"`
	// 	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// 	Name      string    `gorm:"index" json:"name,omitempty"`
	// 	Image     string    `gorm:"uniqueIndex" json:"image,omitempty"`
	// 	FileName  string    `gorm:"uniqueIndex" json:"filename,omitempty"`
	// 	Parent    string    `gorm:"index" json:"parent,omitempty"`
	// 	URL       string    `json:"url,omitempty"`
	// 	Hidden    bool      `json:"hidden,omitempty"`
	// 	Tags      string    `json:"tags,omitempty"`
	// }

	//WAIT FOR PRIMITIVE

	//once primitive returns
	//upload the new image to the database
	//add to the config, the id of the image, and the id of the child/

	//file is now written and be be returned

	marshalledResponse, err := json.Marshal(newModel)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))

}
