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
	// debug := true
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		HandleError(w, errors.New("Missing Parameter UID"))
		return
	}

	//get image and metadata
	model, err := database.Instance.GetImage(uid)
	if err != nil {
		HandleError(w, err)
		return
	}
	tags := r.Header.Get("tags")
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
	config.Input(model.FileName)
	//where do we want to solve the resulting image?
	// i want to make it parentuid-name-file
	// so check incoming filename for path in os...
	// newPath := database.ImageRoot + fmt.Sprintf("%d", model.ParentID)
	// newPath := database.ImageRoot + "temp-" + model.FileName
	//so new image should be saved at /images/parent/image
	//in this case parent = model we currently have

	//how do we want to format outfile?
	output := fmt.Sprintf("%s%d/%s", database.ImageRoot, model.UID, config.Name)
	l.Debug("New image{s} being created at ", output)
	config.Output(output)

	// config.Output = database.ImageRoot + config.Name + fmt.Sprintf("parent-%d", model.UID) + ""

	// //Image will be duplicated to storage...
	// config.output = newPath
	// //start building new image

	//heres what i know, i have a model and an image. And primitive is looking to take that image and write it tothe output
	commandByteArray, _, err := refract.Primitive(*config)
	l.Debug(string(commandByteArray), err)
	if err != nil {
		HandleError(w, err)
		return
	}
	newImages := make([]*database.Model, 0)
	//images are now in their respective parents folder titled
	for i := range config.Outputs {
		name := fmt.Sprintf("%s-%s-%d", config.Name, model.Name, i)
		// output := fmt.Sprintf("%s%d/%s-temp", database.ImageRoot, model.UID, config.Name)

		// image, err := database.Instance.LoadImage(output, 0)
		img, err := ioutil.ReadFile(config.Outputs[i].Path)
		if err != nil {
			HandleError(w, err)
			return
		}
		newHash := md5.Sum(img)
		newFilename := config.Outputs[i].Path
		configString, err := json.Marshal(config)
		if err != nil {
			HandleError(w, err)
			return
		}

		newModel := &database.Model{
			Name:     name,
			Image:    img,
			FileName: newFilename,
			ParentID: model.UID,
			Tags:     fmt.Sprintf("%s %s %s", tags, config.Name, config.Outputs[i].Format),
			FileHash: newHash[:],
			Config:   string(configString),
		}

		//save resulting image
		err = database.Instance.Upload(newModel)
		if err != nil {
			HandleError(w, err)
			return
		}

		//add model to result set
		newImages = append(newImages, newModel)

	}

	marshalledResponse, err := json.Marshal(newImages)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))

}
