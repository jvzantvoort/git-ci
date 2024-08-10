package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/jvzantvoort/git-ci/commands"
	"github.com/jvzantvoort/git-ci/messages"
	log "github.com/sirupsen/logrus"
)

var (
	verbose     bool
	message     string
	scope       string
	programname string
	mtype       string
	subcmnd     string
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
	flag.StringVarP(&mtype, "type", "t", "", "message type (see help)")
	flag.StringVarP(&message, "message", "m", "", "commit message")
	flag.StringVarP(&scope, "scope", "s", "", "type scope (e.g. auth, build, docs, API)")
	flag.ErrHelp = nil
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This application demonstrates extended usage info.\n\n")
		fmt.Fprintf(os.Stderr, "Available flags:\n")
		flag.PrintDefaults()
		PrintTypes()

		os.Exit(0)
	}
}

func main() {

	// Get own name and (git) sub command
	programname = filepath.Base(os.Args[0])
	subcmnd = strings.Replace(programname, "git-", "", 1)

	// Parse the arguments
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if subcmnd == "bs" {
		rand.Seed(time.Now().UnixNano())
		bullshit := messages.GetBullShit()
		randomIndex := rand.Intn(len(bullshit))
		message = bullshit[randomIndex]
		subcmnd = "minor"
	}

	if len(message) == 0 {
		log.Fatalf("Error: missing message")

	}
	args := flag.Args()

	// ci is a special case
	if subcmnd == "ci" {
		subcmnd = "feat"
	}

	if !IsZeroOfUnderlyingType(mtype) {
		subcmnd = mtype
	}

	if subcmnd == "new" {
		subcmnd = "feat"
	}

	msg := NewMessage(message, subcmnd)

	git := commands.NewGitCmd()
	branch := git.Branch()

	err := msg.SetBranch(branch)
	if err != nil {
		log.Warningf("Warning: set branch failed: %s", err)
	}
	msg.SetScope(scope)

	message = msg.String()

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
