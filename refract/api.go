package refract

import (
	"errors"
	"fmt"
	"os/exec"
)

// type API interface {
// 	Primitive(c Config) ([]byte, error)
// }

func Primitive(c Config) ([]byte, error) {
	valid := c.Verify()
	if !valid {
		return nil, errors.New("Invalid config. Aborting")
	}
	out, err := exec.Command(fmt.Sprintf("primitive %s", c.CommandForm())).Output()
	if err != nil {
		return nil, err
	}
	//image will now be in the path specified in config
	return out, nil
}
