package refract

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/happierall/l"
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
	com := fmt.Sprintf("primitive %v", c.CommandForm())
	l.Debug(com)
	out, err := exec.Command("primitive", c.CommandForm()...).Output()
	if err != nil {
		return nil, "", err
	}
	return out, "", nil
}
