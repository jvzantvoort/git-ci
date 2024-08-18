package path

import (
	"errors"
	"os"
	"path"
	"strings"
)

func Which(command string) (string, error) {
	for _, dirname := range strings.Split(os.Getenv("PATH"), ":") {
		fpath := path.Join(dirname, command)
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			continue
		} else {
			return fpath, nil
		}
	}
	return command, errors.New("unable to find command " + command)
}
