package server

import (
	"encoding/json"
	"net/http"
	"refract/database"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func randomImage(w http.ResponseWriter, req *http.Request) {
	model, err := database.Instance.Random()
	if err != nil {
		HandleError(w, err)
		return
	}
	marshalledResponse, err := json.Marshal(model)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}

func getChildren(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		HandleError(w, errors.New("Invalid UID"))
		return
	}
	children, err := database.Instance.GetChildren(uid)
	if err != nil {
		HandleError(w, err)
		return
	}

	marshalledResponse, err := json.Marshal(children)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}

func getImage(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		HandleError(w, errors.New("Invalid UID"))
		return
	}
	model, err := database.Instance.GetImage(uid)
	if err != nil {
		HandleError(w, err)
		return
	}
	marshalledResponse, err := json.Marshal(model)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}

// func getDBStats(w http.ResponseWriter, req *http.Request) {
// 	uid := chi.URLParam(req, "uid")
// 	if uid == "" {
// 		HandleError(w, errors.New("Invalid UID"))
// 		return
// 	}
// 	model, err := database.Instance.GetImage(uid)
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 	marshalledResponse, err := json.Marshal(model)
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 	w.Header().Add("Content-Type", "application/json")
// 	w.Write([]byte(marshalledResponse))
// }
