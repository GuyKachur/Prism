package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"refract/database"
	"strconv"

	"github.com/happierall/l"

	"github.com/pkg/errors"
)

func createImage(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(64 << 20)
	if err != nil {
		HandleError(w, err)
		return
	}
	file, fh, err := req.FormFile("file")
	// fmt.Println("Uploading... ", fh.Filename)
	if err != nil {
		HandleError(w, err)
		return
	}
	defer file.Close()

	md := req.FormValue("input")
	input := &Input{}
	err = json.Unmarshal([]byte(md), input)
	if err != nil {
		HandleError(w, err)
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		HandleError(w, err)
		return
	}
	// img, _, err := image.Decode(bytes.NewReader(b)) //discard options
	// if err != nil {
	// 	HandleError(w, err)
	// 	return
	// }

	parent64, err := strconv.ParseUint(input.Parent, 10, 64)
	parent := uint(parent64)
	if err != nil {
		l.Error(err)
		parent = uint(0)
	}
	fileHash := md5.Sum(b)
	ext := filepath.Ext(fh.Filename)
	fileName := hex.EncodeToString(fileHash[:]) + ext
	// l.Warn(len(fileHash))
	model := &database.Model{
		Name:     input.Name,
		Image:    b,
		FileName: fileName,
		ParentID: parent,
		Tags:     input.Tags,
		FileHash: fileHash[:],
	}
	// fmt.Println("Input: ", input)
	// model.Image = nil
	// fmt.Println("Model: ", model)
	// model.Image = bytes

	err = database.Instance.Upload(model)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`Upload successful: ` + fmt.Sprint(model.UID)))
}

func uploadURL(w http.ResponseWriter, req *http.Request) {
	// url := req.FormValue("url")
	// if url == "" {
	// 	HandleError(w, errors.New("Missing URL in upload body"))
	// }
	md := req.FormValue("metadata")
	input := &Input{}
	err := json.Unmarshal([]byte(md), input)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp, err := OutboundClient.Get(input.URL)
	if err != nil {
		HandleError(w, errors.Wrap(err, "Error downloading image from URL: "))
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		HandleError(w, err)
		return
	}
	//we need to see what the image looks like here, but
	// just make it fileheader

	// img, _, err := image.Decode(bytes.NewReader(b)) //discard options
	// if err != nil {
	// 	HandleError(w, err)
	// 	return
	// }
	parent, err := strconv.ParseUint(input.Parent, 10, 64)
	if err != nil {
		l.Error(err)
	}
	fileHash := md5.Sum(b)
	ext := filepath.Ext(input.URL)
	fileName := hex.EncodeToString(fileHash[:]) + ext
	// bytes.ToValidUTF8(b, make([]byte, 1))
	model := &database.Model{
		Name:     input.Name,
		Image:    b,
		FileName: fileName,
		ParentID: uint(parent),
		URL:      input.URL,
		Tags:     input.Tags,
		FileHash: fileHash[:],
	}

	err = database.Instance.Upload(model)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`URL Upload successful:` + fmt.Sprint(model.UID)))
}
