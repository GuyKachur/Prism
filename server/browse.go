package server

import (
	"encoding/json"
	"net/http"
	"refract/database"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func browse(w http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(req, "page"))
	if err != nil {
		handleError(w, errors.New("Invalid page"))
		return
	}
	pageSize, err := strconv.Atoi(chi.URLParam(req, "pageSize"))
	if err != nil {
		handleError(w, errors.New("Invalid pageSize"))
		return
	}
	images, err := database.Instance.LoadImages(page, pageSize)
	if err != nil {
		handleError(w, err)
		return
	}

	marshalledResponse, err := json.Marshal(images)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}
