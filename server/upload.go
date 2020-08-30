package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"refract/database"

	"github.com/pkg/errors"
)

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
		UID:      0,
		Name:     input.Name,
		Image:    bytes,
		FileName: input.Extension,
		Parent:   input.Parent,
		Tags:     input.Tags,
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

func uploadURL(w http.ResponseWriter, req *http.Request) {
	url := req.FormValue("url")
	if url == "" {
		handleError(w, errors.New("Missing URL in upload body"))
	}
	md := req.FormValue("metadata")
	input := &Input{}
	err := json.Unmarshal([]byte(md), input)
	if err != nil {
		handleError(w, err)
		return
	}

	resp, err := OutboundClient.Get(url)
	if err != nil {
		handleError(w, errors.Wrap(err, "Error downloading image from URL: "))
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	model := &database.Model{
		UID:      0,
		Name:     input.Name,
		Image:    bytes,
		FileName: input.Extension,
		Parent:   input.Parent,
		URL:      url,
		Tags:     input.Tags,
	}

	err = database.Instance.Upload(model)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`URL Upload successful!`))
}
