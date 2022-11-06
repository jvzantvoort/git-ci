package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

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

	flag.BoolVar(&verbose, "verbose", false, "verbose messages")
	flag.BoolVar(&verbose, "v", false, "verbose messages")
	flag.StringVar(&message, "m", "", "commit message")
	flag.StringVar(&message, "message", "", "commit message")

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

	match := pattern.FindStringSubmatch(branch)

	if len(match) != 0 {
		lastIndex := pattern.SubexpIndex("ticket")
		ticket := match[lastIndex]
		message = fmt.Sprintf("%s %s", ticket, message)
	}
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
