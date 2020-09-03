package server

import (
	"encoding/json"
	"net/http"
	"refract/database"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func search(w http.ResponseWriter, req *http.Request) {
	term := chi.URLParam(req, "term")
	if term == "" {
		HandleError(w, errors.New("Invalid term"))
		return
	}
	images, err := database.Instance.Search(term)
	if err != nil {
		HandleError(w, err)
		return
	}

	marshalledResponse, err := json.Marshal(images)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}
