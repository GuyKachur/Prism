package server

import (
	"net/http"
	"refract/database"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type Input struct {
	Name   string
	Parent string
	Tags   string
}

func artHandler() http.Handler {
	//register all the handlers
	r := chi.NewRouter()
	r.Get("/image/{uid}", getImage)
	r.Get("/image/{uid}/children", getChildren)
	r.Get("/image/random", randomImage)
	r.Delete("/delete/{uid}", deleteImage)
	r.Post("/upload", createImage)
	r.Post("/upload/url", uploadURL)
	r.Get("/browse/{page}/{pageSize}", Browse)
	r.Get("/search/{term}", search)
	return r
}

func deleteImage(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		HandleError(w, errors.New("Invalid UID"))
		return
	}
	err := database.Instance.Delete(uid)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`Delete Successful`))
}
