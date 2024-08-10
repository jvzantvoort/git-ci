package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/jvzantvoort/git-ci/commands"
	log "github.com/sirupsen/logrus"
)

var (
	verbose bool
	message string
	pattern = regexp.MustCompile(`^(?P<type>\w+)/(?P<ticket>\w+-\d+).*$`)
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	flag.BoolVarP(&verbose, "verbose", "v", false, "verbose messages")
	flag.StringVarP(&message, "message", "m", "", "commit message")

}

func ExtractTicket(instr string) (string, error) {
	retv := ""
	var err error
	match := pattern.FindStringSubmatch(instr)

	if len(match) != 0 {
		lastIndex := pattern.SubexpIndex("ticket")
		retv = match[lastIndex]
		err = nil
	} else {
		err = fmt.Errorf("cannot find ticket in string")
	}

	return retv, err
}

func BuildMessageString(branchname, message string) (string, error) {
	if len(message) == 0 {
		return "", fmt.Errorf("message is empty")
	}
	if len(branchname) == 0 {
		return "", fmt.Errorf("branch is empty")
	}

	retv := message

	ticket, err := ExtractTicket(branchname)
	if err == nil {
		retv = fmt.Sprintf("%s %s", ticket, message)
	}
	return retv, nil
}

func main() {
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if len(message) == 0 {
		log.Fatalf("Error: missing message")

	}
	args := flag.Args()

	log.Debugf("actions: %s\n", strings.Join(args, " "))

	git := commands.NewGitCmd()
	branch := git.Branch()

	log.Debugf("branch:  %s\n", branch)

	message, err := BuildMessageString(branch, message)
	if err != nil {
		log.Errorf("message failed: %s", err)
		return
	}

	log.Debugf("message: %s", message)

	stdout_lines, stderr_lines, exitcode := git.Commit(message, args...)

	if exitcode != nil {
		log.Errorf("Exit code: %s", exitcode)
	}
	for _, line := range stdout_lines {
		fmt.Println(line)
	}
	for _, line := range stderr_lines {
		fmt.Printf("err: %s\n", line)
	}

}
