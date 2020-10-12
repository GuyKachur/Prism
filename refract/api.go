package refract

import (
	"errors"
	"os/exec"
)

// type API interface {
// 	Primitive(c Config) ([]byte, error)
// }

func Primitive(c Config) ([]byte, string, error) {
	//This needs to have a deterministic, but concurrency safe way of saving new images...

	valid := c.Verify()
	if !valid {
		return nil, "", errors.New("Invalid config. Aborting")
	}
	out, err := exec.Command("primitive", c.CommandForm()...).Output()
	if err != nil {
		return nil, "", err
	}
	return out, "", nil
}
