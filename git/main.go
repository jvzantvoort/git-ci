package git

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"

	"strings"

	"github.com/jvzantvoort/git-ci/utils"
	"github.com/jvzantvoort/git-ci/path"
	log "github.com/sirupsen/logrus"
)

// GitCmd object for git
type GitCmd struct {
	Path            *path.Path
	Cwd             string
	Command         string
	CommandMap      map[string]string
	FormatKeyNames  []string
	FormatKeyValues []string
}

func (g GitCmd) Execute(args ...string) ([]string, []string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	stdout_list := []string{}
	stderr_list := []string{}
	cmnd := []string{}

	cmnd = append(cmnd, args...)

	utils.Debugf("command %s %s", g.Command, strings.Join(cmnd, " "))

	cmd := exec.Command(g.Command, cmnd...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("stdout pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Errorf("stderr pipe failed, %v", err)
		log.Fatal(err)
		panic(err)
	}

	cmd.Dir = g.Cwd
	err = cmd.Start()
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	stdout_scan := bufio.NewScanner(stdout)
	stdout_scan.Split(bufio.ScanLines)
	for stdout_scan.Scan() {
		msg := stdout_scan.Text()
		stdout_list = append(stdout_list, msg)
	}

	stderr_scan := bufio.NewScanner(stderr)
	stderr_scan.Split(bufio.ScanLines)
	for stderr_scan.Scan() {
		msg := stderr_scan.Text()
		stderr_list = append(stderr_list, msg)
	}

	eerror := cmd.Wait()
	if eerror != nil {
		log.Errorf("command failed, %v", eerror)
	}
	return stdout_list, stderr_list, eerror
}

func (g GitCmd) ExecuteSingle(args ...string) string {
	stdout, _, _ := g.Execute(args...)
	if len(stdout) == 0 {
		return ""
	}
	return string(stdout[0])
}

// Aliasses
// URL function returning the git url
func (g GitCmd) URL() string {
	return g.ExecuteSingle("config", "--get", "remote.origin.url")
}

// Branch function returning the current git branch
func (g GitCmd) Branch() string {
	return g.ExecuteSingle("rev-parse", "--abbrev-ref", "HEAD")
}

// Root function returning the git root
func (g GitCmd) Root() string {
	return g.ExecuteSingle("rev-parse", "--show-toplevel")
}

// New create a new git object
func New() *GitCmd {
	utils.LogStart()
	defer utils.LogEnd()

	retv := &GitCmd{}
	retv.Initialize()
	retv.Path = path.New("PATH")

	dir, err := os.Getwd()
	if err != nil {
		utils.Fatalf("%s", err)
	} else {
		retv.Cwd = dir
	}

	if result, err := retv.Path.LookupPlatform(retv.CommandMap); err == nil {
		retv.Command = result
	}

	return retv
}

func (g GitCmd) Commit(message string, args ...string) ([]string, []string, error) {
	utils.LogStart()
	defer utils.LogEnd()

	arglist := []string{}
	arglist = append(arglist, "commit")
	arglist = append(arglist, "--message")
	arglist = append(arglist, message)
	arglist = append(arglist, args...)

	return g.Execute(arglist...)
}

func (g GitCmd) Log(args ...string) (string, error) {

	var stdout bytes.Buffer

	arguments := []string{"log"}

	arguments = append(arguments, args...)

	cmd := exec.Command(g.Command, arguments...)

	cmd.Stdout = &stdout
	cmd.Dir = g.Cwd

	err := cmd.Run()
	outStr := stdout.String()

	return outStr, err
}
