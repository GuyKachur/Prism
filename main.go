package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"refract/server"

	"github.com/happierall/l"
	"github.com/pkg/errors"
)

//modify power level color theme -
const (
	colorHolder  = "\033[38;5;%dm%s\033[39;49m "
	errorColor   = 203
	successColor = 49
)

// func PrintError(msg string) {
// 	PrintColor(errorColor, msg)
// }

func PrintColor(color int, msg string) {
	fmt.Printf(colorHolder, color, msg)
	fmt.Println()
}

type Config struct {
	InputFilePath string `json:"input_file_path,omitempty"`
	LoadFiles     bool   `json:"load_files,omitempty"`
	Scribe        bool   `json:"scribe,omitempty"`
	Force         bool   `json:"force,omitempty"`
}

var Seen map[string]bool

// var Seen map[string]string
// var files []string

func runConfig(config Config) {
	defer func() {
		if r := recover(); r != nil {
			l.Error("Config failed to load: %v", r)
			if config.Scribe {
				l.Debug("The system blinks off for a moment before returning to life... Systems operational")

			}
		}
	}()
	if config.LoadFiles {
		if config.InputFilePath != "" {
			if config.Scribe {
				l.Debug("The screen above lights up again, deBUGing #%^#$... beeps start coming out of a panel below as a small fan whirs to life beneath you.")
			}
			err := filepath.Walk(config.InputFilePath, func(path string, info os.FileInfo, err error) error {
				l.Debug(path)
				if path != "" {
					if !Seen[path] {
						if ok := upload(path); ok {
							Seen[path] = true
						}
						l.Debugf("Skipping: %s", path)

					}
				}

				return nil
			})
			if err != nil {
				panic(err)
			}
			err = save()
			if err != nil {
				panic(err)
			}

		}
	}
	//we arent loading the files what are we doing?
}
func save() error {
	seenStr, err := json.Marshal(Seen)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./images/seen.json", seenStr, 0764)
	if err != nil {
		return err
	}
	return nil
}

func load(fn string) error {
	file, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &Seen)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("File empty")
	}
	return nil
}

func upload(fn string) bool {
	file, err := os.Open(fn)
	if err != nil {
		l.Error(err)
		return false
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		l.Error(err)
		return false
	}

	fi, err := file.Stat()
	if err != nil {
		l.Error(err)
		return false
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		l.Error(err)
		return false
	}
	_, err = part.Write(fileContents)
	if err != nil {
		l.Error(err)
		return false
	}

	input, err := json.Marshal(server.Input{
		Name: fn,
		Tags: "toGeo",
	})
	if err != nil {
		l.Error(err)
		return false
	}

	err = writer.WriteField("input", string(input))
	if err != nil {
		l.Error(err)
		return false
	}
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:9090/upload", body)
	if err != nil {
		l.Error(err)
		return false
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		l.Error(err)
		return false
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	l.Debug(string(respBody))

	return true
}

func init() {
	//setup seen map
	//try and load, if not make

	if err := load("./images/seen.json"); err != nil {
		l.Error(err)
	}

}

func main() {
	defer func() {
		if r := recover(); r != nil {
			//attempting to save seen to file
			if err := save(); err != nil {
				l.Error(err)
				return
			}
			l.Error("and suddenly, the cacophony of noises that had filled the room abruptly end. the screen giving one flash of bright before shutting off. The room plunges into darkness.")
			l.Log("Save successful")
		}
	}()

	if configBytes, err := ioutil.ReadFile("config.json"); err != nil {
		l.Error(err)
		panic(err)
	} else {
		config := Config{}
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			l.Error(err)
		}
		go runConfig(config)
	}

	// envConfig := os.Getenv("PRISM_CONFIG")
	// config := Config{}

	// err := json.Unmarshal([]byte(envConfig), &config)
	// if err == nil {
	// 	//TODO: Load scene
	// 	go runConfig(config)
	// }
	PrintColor(successColor, "-_-Starting Server-_-")
	PrintColor(successColor, "A small LED screen lights up... [:9090] blinks merrily back at you")
	server.NewServer()
	// for i := 1; i < 10; i++ {
	// 	fout, err := exec.Command("primitive", "-i", "/home/guy/projects/go/geo/prism/images/a.jpg", "-o", fmt.Sprintf("/home/guy/projects/go/geo/prism/images/output/a-shapes-%d.svg", i),
	// 		"-n", fmt.Sprint(i*50)).Output()
	// 	if err != nil {
	// 		l.Error(err)
	// 		return
	// 	}
	// 	l.Debug(fout)
	// }

	// con := refract.CreateDefault()
	// con.Input = "/home/guy/projects/go/geo/prism/images/a.jpg"
	// con.Output = "/home/guy/projects/go/geo/prism/images/a.svg"
	// con.Verbose = "on"
	// out, _, err := refract.Primitive(*con)
	// if err != nil {
	// 	l.Error(err)
	// 	return
	// }
	// PrintColor(successColor,
	// 	fmt.Sprintf("Sucess! config: %s", con.CommandForm()))
	// fmt.Println(string(out))
}

// func main() {
// 	for j := 0; j < 256; j++ {
// 		fmt.Printf(colorHolder, j, fmt.Sprintf("%d", j))
// 		PrintColor(j, "Hello!")
// 	}
// }
