package server

import (
	"encoding/json"
	"net/http"
	"refract/database"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

// type Response struct {
// 	rows []database.Model `json:"rows,omitempty"`
// }

func Browse(w http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(req, "page"))
	if err != nil {
		HandleError(w, errors.New("Invalid page"))
		return
	}
	pageSize, err := strconv.Atoi(chi.URLParam(req, "pageSize"))
	if err != nil {
		HandleError(w, errors.New("Invalid pageSize"))
		return
	}
	images, err := database.Instance.LoadImages(page, pageSize)
	if err != nil {
		HandleError(w, err)
		return
	}
	// pageResponse := Response{
	// 	rows: *images,
	// }

	marshalledResponse, err := json.Marshal(images)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}
