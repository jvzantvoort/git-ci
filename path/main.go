// PATH type handling.
package path

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/jvzantvoort/git-ci/utils"
	"github.com/mitchellh/go-homedir"
)

type Path struct {
	Type        string
	Home        string
	Directories []string
}

func (p Path) HavePath(inputdir string) bool {
	utils.LogStart()
	defer utils.LogEnd()

	for _, element := range p.Directories {
		if element == inputdir {
			return true
		}
	}
	return false
}

// AppendPath append a path to the list of Directories
func (p *Path) AppendPath(inputdir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	utils.Debugf("inputdir=%s", inputdir)

	if len(inputdir) == 0 {
		return nil
	}

	fullpath, err := homedir.Expand(inputdir)
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}

	_, err = os.Stat(fullpath)
	if err != nil {
		return err
	}

	if !p.HavePath(fullpath) {
		p.Directories = append(p.Directories, fullpath)
	}

	return nil
}

func (p *Path) PrependPath(inputdir string) error {
	utils.LogStart()
	defer utils.LogEnd()

	utils.Debugf("inputdir=%s", inputdir)

	if len(inputdir) == 0 {
		return nil
	}

	fullpath, err := homedir.Expand(inputdir)
	if err != nil {
		utils.Errorf("error %s", err)
		return err
	}

	if !p.HavePath(fullpath) {
		p.Directories = append([]string{fullpath}, p.Directories...)
	}

	return nil
}

func (p *Path) Import(path string) {
	utils.LogStart()
	defer utils.LogEnd()

	utils.Debugf("path=%s", path)

	for _, dirn := range strings.Split(path, ":") {
		err := p.AppendPath(dirn)
		if err != nil {
			utils.Errorf("Error: %v", err)
		}
	}
}

func (p Path) IsEmpty() bool {
	if len(p.Directories) == 0 {
		return true
	} else {
		return false
	}
}

func (p Path) ReturnExport() string {
	return fmt.Sprintf("export %s=\"%s\"", p.Type, strings.Join(p.Directories, ":"))

}

// targetExists return true if target exists
func (p Path) targetExists(targetpath string) bool {
	_, err := os.Stat(targetpath)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (p Path) Lookup(target string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	var retv string
	var err error
	err = fmt.Errorf("command %s not found", target)

	for _, dirname := range p.Directories {
		fullpath := path.Join(dirname, target)
		if p.targetExists(fullpath) {
			retv = fullpath
			err = nil
			break
		}
	}
	return retv, err
}

func (p Path) LookupMulti(targets ...string) (string, error) {
	for _, target := range targets {
		if result, err := p.Lookup(target); err != nil {
			return result, nil
		}
	}
	return "", fmt.Errorf("targets not found")
}

func (p Path) MapGetPlatform(pathmap map[string]string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	goos := runtime.GOOS
	utils.Debugf("os=%s", goos)
	if target, ok := pathmap[goos]; ok {
		utils.Debugf("found key: %s -> %s", goos, target)
		return target, nil
	}
	if target, ok := pathmap["default"]; ok {
		utils.Debugf("found key: default -> %s", target)
		return target, nil
	}

	return "", fmt.Errorf("map keys not found")
}

// LookupPlatform lookup paths based on platform
func (p Path) LookupPlatform(pathmap map[string]string) (string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	commandname, err := p.MapGetPlatform(pathmap)
	if err != nil {
		utils.Errorf("Err: %s", err)
		return "", err
	}

	if result, err := p.Lookup(commandname); err == nil {
		utils.Debugf("found: %s -> %s", commandname, result)
		return result, nil
	}
	utils.Errorf("cannot find %s in path", commandname)

	return "", fmt.Errorf("target not found")
}

func (p Path) Git() string {

	CommandMap := map[string]string{
		"windows": "git.exe",
		"linux":   "git",
		"default": "git",
	}

	if result, err := p.LookupPlatform(CommandMap); err == nil {
		return result
	}
	return ""
}

func New(pathname string) *Path {
	retv := &Path{}
	retv.Type = pathname
	retv.Home, _ = homedir.Dir()
	retv.Import(os.Getenv(retv.Type))
	return retv
}
