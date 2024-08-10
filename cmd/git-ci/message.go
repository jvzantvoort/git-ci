package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Message struct {
	Type    string
	Scope   string
	Ticket  string
	Message string
	Branch  string
}

var (
	pattern      = regexp.MustCompile(`^(?P<type>\w+)/(?P<ticket>\w+-\d+).*$`)
	ValidSubSets = []string{
		"ci",
		"feat",
		"fix",
		"chore",
		"docs",
		"style",
		"refactor",
		"perf",
		"test",
		"new",
		"build",
		"revert",
		"config",
		"hotfix",
		"release",
	}
	commitTypes = map[string]string{
		"feat":     "A new feature for the user.",
		"fix":      "A bug fix.",
		"chore":    "Routine task that doesn’t modify application logic.",
		"docs":     "Documentation-only changes.",
		"style":    "Changes that do not affect the meaning of the code (white-space, formatting, etc.).",
		"refactor": "A code change that neither fixes a bug nor adds a feature but improves the code.",
		"perf":     "A code change that improves performance.",
		"test":     "Adding missing tests or correcting existing tests.",
		"build":    "Changes that affect the build system or external dependencies.",
		"revert":   "Reverts a previous commit.",
		"config":   "Changes to project configuration files.",
		"hotfix":   "A critical fix that needs to be made immediately.",
		"release":  "A commit related to releasing a new version or a release process.",
	}
)

func PrintTypes() {
	fmt.Fprintf(os.Stderr, "\n\nFor conveniance the command can be aliased to:\n")
	for _, element := range ValidSubSets {
		fmt.Fprintf(os.Stderr, "  git-%s\n", element)
	}
	fmt.Fprintf(os.Stderr, "\n\nExample:\n\n")
	fmt.Fprintf(os.Stderr, "    git feat -s api -m \"message\"\n\n")
	fmt.Fprintf(os.Stderr, "    git new -s api -m \"Added foo interface\"\n\n")
	fmt.Fprintf(os.Stderr, "\n\nDescription of types:\n\n")
	for name, descr := range commitTypes {
		fmt.Fprintf(os.Stderr, "  %-12s %s\n", name, descr)
	}
	fmt.Fprintf(os.Stderr, "\n\n")
}

func NewMessage(message, subcmnd string) *Message {
	retv := &Message{}
	retv.Message = message
	if stringInSlice(subcmnd, ValidSubSets) {
		retv.Type = subcmnd
	}

	return retv
}

func (m *Message) SetBranch(branch string) error {
	m.Branch = branch
	if len(m.Branch) == 0 {
		return fmt.Errorf("branch is empty")
	}
	ticket, err := m.ExtractTicket(branch)
	if err != nil {
		return err
	}
	m.Ticket = ticket
	return nil
}

func (m *Message) SetScope(scope string) {
	if !IsZeroOfUnderlyingType(scope) {
		m.Scope = scope
	}
}

func (m Message) ExtractTicket(instr string) (string, error) {
	var retv string
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

func (m Message) String() string {
	var retv string

	if !IsZeroOfUnderlyingType(m.Type) {
		if IsZeroOfUnderlyingType(m.Scope) {
			retv = fmt.Sprintf("%s:", m.Type)
		} else {
			retv = fmt.Sprintf("%s(%s):", m.Type, m.Scope)
		}
	}

	if !IsZeroOfUnderlyingType(m.Ticket) {
		retv = fmt.Sprintf("%s %s", retv, m.Ticket)
	}
	retv = fmt.Sprintf("%s %s", retv, m.Message)
	retv = strings.TrimSpace(retv)
	return retv
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
