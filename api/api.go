package api

import (
	"errors"
	"fmt"
	"os/exec"
)

func Primitive(c Config) ([]byte, error) {
	valid := c.Verify()
	if !valid {
		return nil, errors.New("Invalid config. Aborting")
	}
	out, err := exec.Command(fmt.Sprintf("primitive %s", c.CommandForm())).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
