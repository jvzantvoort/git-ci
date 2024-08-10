package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GitCmd object for git
type GitCmd struct {
	Path       *Path
	Cwd        string
	Command    string
	CommandMap map[string]string
}

func (g GitCmd) Prefix() string {
	pc, _, _, _ := runtime.Caller(1)
	elements := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return fmt.Sprintf("GitCmd.%s", elements[len(elements)-1])
}

func (g GitCmd) Execute(args ...string) ([]string, []string, error) {
	log_prefix := g.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	stdout_list := []string{}
	stderr_list := []string{}
	cmnd := []string{}

	cmnd = append(cmnd, args...)

	log.Debugf("%s: command %s %s", log_prefix, g.Command, strings.Join(cmnd, " "))

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

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		msg := scanner.Text()
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

// NewGitCmd create a new git object
func NewGitCmd() *GitCmd {
	retv := &GitCmd{}

	log_prefix := retv.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	retv.Path = NewPath("PATH")

	retv.CommandMap = map[string]string{
		"windows": "git.exe",
		"linux":   "git",
		"default": "git",
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%s: %s", log_prefix, err)
	} else {
		retv.Cwd = dir
	}

	if result, err := retv.Path.LookupPlatform(retv.CommandMap); err == nil {
		retv.Command = result
	}

	return retv
}

// Aliasses
// URL function returning the git url
func (g GitCmd) URL() string {
	stdout, _, _ := g.Execute("config", "--get", "remote.origin.url")
	if len(stdout) == 0 {
		return ""
	}
	return string(stdout[0])
}

// Branch function returning the current git branch
func (g GitCmd) Branch() string {
	stdout, _, _ := g.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if len(stdout) == 0 {
		return ""
	}
	return string(stdout[0])
}

// Root function returning the git root
func (g GitCmd) Root() string {
	stdout, _, _ := g.Execute("rev-parse", "--show-toplevel")
	if len(stdout) == 0 {
		return ""
	}
	return string(stdout[0])
}

func (g GitCmd) Commit(message string, args ...string) ([]string, []string, error) {
	log_prefix := g.Prefix()
	log.Debugf("%s: start", log_prefix)
	defer log.Debugf("%s: end", log_prefix)

	arglist := []string{}
	arglist = append(arglist, "commit")
	arglist = append(arglist, "--message")
	arglist = append(arglist, message)
	arglist = append(arglist, args...)

	return g.Execute(arglist...)
}
