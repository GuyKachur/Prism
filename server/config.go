package server

import (
	"database"
	"encoding/json"
	"errors"
	"image"
	"net/http"
	"refract/server"

	"github.com/go-chi/chi"
)

type RefractRequest struct {
	image  image.Image
	config api.Config
}

func ConfigHandler() http.Handler {
	//register all the handlers
	r := chi.NewRouter()
	r.Get("/refract/{uid}", refractImage)
	return r
}

/

//Should profiles just exist on the front end? => no probably not right???
//if uid === url? that would let me dictate logic and collapse those handlers. Though tbh too much is done in the handlers anyways
