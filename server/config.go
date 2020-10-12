package server

import (
	"image"
	"net/http"
	"refract/refract"

	"github.com/go-chi/chi"
)

type RefractRequest struct {
	image  image.Image
	config refract.Config
}

func ConfigHandler() http.Handler {
	//register all the handlers
	r := chi.NewRouter()
	r.Get("/refract/{uid}", RefractHandler)
	return r
}

//Should profiles just exist on the front end? => no probably not right???
//if uid === url? that would let me dictate logic and collapse those handlers. Though tbh too much is done in the handlers anyways
