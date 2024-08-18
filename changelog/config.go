package changelog

import (
	"path/filepath"
	"strings"

	"github.com/jvzantvoort/git-ci/git"
)

const (
	configFile       string = ".changelog.yml"
	commentDelimiter string = "%x00"
)

type Changelog struct {
	Command         git.GitCmd
	FormatKeyNames  []string
	FormatKeyValues []string
}

func (c *Changelog) Initialize() {
	for keyname, value := range GitFormatKeys {
		c.FormatKeyNames = append(c.FormatKeyNames, keyname)
		c.FormatKeyValues = append(c.FormatKeyValues, value)
	}
}

func New() *Changelog {
	retv := &Changelog{}
	retv.Command = git.New()
	retv.Initialize()
	return retv
}

func (c Changelog) Configfile() string {
	path := c.Command.Root()
	path = filepath.Join(path, configFile)
	return path
}

func (c Changelog) Command() {
	format := strings.Join(c.FormatKeyValues, commentDelimiter)
	args = []string{}
	args = append(args, "-z")
	args = append(args, "--topo-order")
	args = append(args, "--pretty=format:"+format)

	stdout_list, stderr_list, eerror := c.Command.Log(args...)

}
