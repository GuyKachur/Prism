package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"refract/database"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func artHandler() http.Handler {
	//register all the handlers
	r := chi.NewRouter()
	r.Get("/image/{uid}", getImage)
	r.Get("/image/{uid}/children", getChildren)
	r.Delete("/image/{uid}", deleteImage)
	r.Post("/create", createImage)
	r.Get("/browse/{page}/{pageSize}", listImages)

	return r
}

func listImages(w http.ResponseWriter, req *http.Request) {
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

func deleteImage(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		handleError(w, errors.New("Invalid UID"))
		return
	}
	err := database.Instance.Delete(uid)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`Delete Successful`))
}

type Input struct {
	Name      string
	Extension string
	Parent    string
}

func createImage(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(64 << 20)
	if err != nil {
		handleError(w, err)
		return
	}
	file, fh, err := req.FormFile("input")
	fmt.Println("Uploading... ", fh.Filename)
	if err != nil {
		handleError(w, err)
		return
	}
	defer file.Close()

	md := req.FormValue("metadata")
	input := &Input{}
	err = json.Unmarshal([]byte(md), input)
	if err != nil {
		handleError(w, err)
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(w, err)
		return
	}

	model := &database.Model{
		UID:       0,
		Name:      input.Name,
		Image:     bytes,
		Extension: input.Extension,
		Parent:    input.Parent,
	}
	// fmt.Println("Input: ", input)
	// model.Image = nil
	// fmt.Println("Model: ", model)
	// model.Image = bytes

	err = database.Instance.Upload(model)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`Upload successful!`))
}

func getChildren(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		handleError(w, errors.New("Invalid UID"))
		return
	}
	children, err := database.Instance.GetChildren(uid)
	if err != nil {
		handleError(w, err)
		return
	}

	marshalledResponse, err := json.Marshal(children)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}

func getImage(w http.ResponseWriter, req *http.Request) {
	uid := chi.URLParam(req, "uid")
	if uid == "" {
		handleError(w, errors.New("Invalid UID"))
		return
	}
	model, err := database.Instance.GetImage(uid)
	if err != nil {
		handleError(w, err)
		return
	}
	marshalledResponse, err := json.Marshal(model)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(marshalledResponse))
}

func handleError(w http.ResponseWriter, err error) {
	//LOG???
	http.Error(w, err.Error(), 500)
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message":"pong"}`))
}

func MainServer() {
	r := chi.NewRouter()
	//gochi suggested middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", healthCheck)
	r.Mount("/art", artHandler())
	http.ListenAndServe(":9090", r)
}

// art/get/{uid}
// art/get/{uid}/GetChildren
// art/create/
// // type Datastore interface {
// 	GetItem(uid string) (*Model, error)
// 	GetChildren(parentUID string) (*[]Model, error)
// 	Upload(model *Model) error
// 	// Update(model *Model) error
// 	Delete(model *Model) error
// 	LoadImages(page, pageSize int) (*[]Model, error)
// }
// ##Routes
// /art/generate/
// /art/get/{uid}
// /profile/register
// 	-this could just always update edit
