package refract

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

//Config represents the options that primitive accepts
//version config?
type Config struct {
	gorm.Model
	input       string
	output      string
	Name        string   `gorm:"unique" json:"name,omitempty"`
	Number      int      `json:"number,omitempty"`
	Mode        int      `json:"mode,omitempty"`        //0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect, 6=beziers, 7=rotatedellipse, 8=polygon
	Rep         int      `json:"rep,omitempty"`         //add N extra shapes each iteration with reduced search (mostly good for beziers)
	Nth         int      `json:"nth,omitempty"`         //save every nth frame, %d must be in output Path
	InSize      int      `json:"in_size,omitempty"`     //size to resize large input images before processing
	OutSize     int      `json:"out_size,omitempty"`    //output image size
	Alpha       int      `json:"alpha,omitempty"`       //color alpha, 0=algos choice
	Background  string   `json:"background,omitempty"`  //avg, 	starting background color (hex)
	Workers     int      `json:"workers,omitempty"`     // number of parallel workers, 0 = all cores
	Verbose     string   `json:"verbose,omitempty"`     //verbose: off
	VeryVerbose string   `json:"veryverbose,omitempty"` //very verbose: off
	Outputs     []Output `gorm:"-" json:"outputs,omitempty"`
}

type Output struct {
	gorm.Model
	// ID   string `gorm:"primary
	Format string
	Path   string // png jpg svg gif
}

// func (of *Output) GetType() string {
// 	return path.Ext(of.Path
// }

// {"name":"Unknown","number":"35","mode":"0","rep":"0","nth":"1","in_size":"256","out_size":"512","alpha":"0","workers":"0","verbose":"on","veryverbbose":"off","Outputs":[{"format":".svg"}]}

//CreateDefault creats a default config, missing only input and output
func CreateDefault() *Config {
	outputs := make([]Output, 1)
	outputs[0].Format = ".jpg"
	return &Config{
		input:       "",
		output:      "",
		Number:      100,
		Mode:        0,
		Rep:         0,
		Nth:         1,
		InSize:      256,
		OutSize:     1024,
		Alpha:       0,
		Background:  "avg",
		Workers:     0,
		Verbose:     "on",
		VeryVerbose: "off",
		Outputs:     outputs,
	}
}

//Verify returns true if config is runnable
func (c *Config) Verify() bool {
	if c.Outputs == nil {
		return false
	}
	// for i := range c.Outputs {
	// 	ext := Path.Ext(c.Outputs[i].filename)
	// 	switch ext {
	// 	case ".svg":
	// 		{
	// 			c.Outputs[i].Format = "svg"
	// 		}
	// 	}

	// }
	return c.input != "" && c.output != "" && c.Number != 0
}

//CommandForm returns the config in the form expected by the 'primitive' command
func (c *Config) CommandForm() []string {
	outputString := ""
	//error check!
	for i := range c.Outputs {
		outputString = outputString + fmt.Sprintf("-o %s", c.Outputs[i].Path)
	}

	return strings.Split(fmt.Sprintf("-i %s %s -n %d -m %d -rep %d -nth %d -r %d -s %d -a %d -bg %s -j %d -v %s -vv %s",
		c.input,
		outputString,
		c.Number,
		c.Mode,
		c.Rep,
		c.Nth,
		c.InSize,
		c.OutSize,
		c.Alpha,
		c.Background,
		c.Workers,
		c.Verbose,
		c.VeryVerbose), " ")
}
func (c *Config) Input(inPath string) {
	c.input = "./images/" + inPath
}
func (c *Config) Output(outPath string) {
	c.output = outPath
	for i := range c.Outputs {
		c.Outputs[i].Path = fmt.Sprintf("%s-%d%s", outPath, time.Now().UnixNano(), c.Outputs[i].Format)
	}
}

// fmt.Println("YYYY-MM-DD hh:mm:ss : ", currentTime.Format("2006-01-02 15:04:05"))
