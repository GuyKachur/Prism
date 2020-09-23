package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"refract/database"
	"refract/server"

	"github.com/go-chi/chi"
	"github.com/happierall/l"
	"github.com/pkg/errors"
)

//assumes Image
func RefractHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid != "" {
		server.HandleError(w, errors.New("Missing Parameter UID"))
	}
	//get image and metadata
	model, err := database.Instance.GetImage(uid)
	if err != nil {
		server.HandleError(w, err)
		return
	}

	ca := r.Header.Get("config")
	config := &api.Config{}
	err = json.Unmarshal([]byte(ca), config)
	if err != nil {
		server.HandleError(w, err)
		return
	}
	// save image and spit path out into config
	path, err := database.Instance.SaveImage()
	if err != nil {
		server.HandleError(w, err)
		return
	}
	config.input = path
	//Image will be duplicated to storage...
	config.output = imageRoot + fmt.Sprintf("%d", model.UID) + "-" + model.FileName
	commandByteArray, err = api.Primitive(*config)
	l.Debug(commandByteArray)
	if err != nil {
		server.HandleError(w, err)
		return
	}
	name := r.Header.Get("name")
	if name != "" {
		name = config.Name + "-" + model.Name
	}
	
	//that will hopefully return a success.... though it might only return on error state
	//fill the output path
	//now we need to grab the outputted image,
	newImage := &Model{
		Name: name
	}

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

	//once primitive returns
	//upload the new image to the database
	//add to the config, the id of the image, and the id of the child/

	//file is now written and be be returned

	marshalledResponse, err := json.Marshal(model)
	if err != nil {
		server.HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))

}
