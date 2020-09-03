package profile

import (
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"refract/api"
	"refract/database"
	"refract/server"

	"github.com/go-chi/chi"
)

type RefractRequest struct {
	image  image.Image
	config api.Config
}

func ProfileHandler() http.Handler {
	//register all the handlers
	r := chi.NewRouter()
	r.Get("/refract/{uid}", refractImage)
	return r
}

//assumes Image
func refractImage(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid != "" {
		server.HandleError(w, errors.New("Missing Parameter UID"))
	}
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

	newImage, err = api.Primitive(*config)
	if err != nil {
		server.HandleError(w, err)
		return
	}
	//file is now written and be be returned

	marshalledResponse, err := json.Marshal(model)
	if err != nil {
		server.HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))

}

func primitive() {

}

//Should profiles just exist on the front end? => no probably not right???
//if uid === url? that would let me dictate logic and collapse those handlers. Though tbh too much is done in the handlers anyways
