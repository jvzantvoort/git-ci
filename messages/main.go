package messages

import (
	"embed"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Content missing godoc.
//
//go:embed bullshit/*
var Content embed.FS

func GetContent(folder, name string) string {
	filename := fmt.Sprintf("%s/%s", folder, name)

	msgstr, err := Content.ReadFile(filename)
	if err != nil {
		log.Errorf("%s", err)
		msgstr = []byte("undefined")
	}
	return strings.TrimSuffix(string(msgstr), "\n")

}

func GetBullShit() []string {
	retv := []string{}
	content := GetContent("bullshit", "commit")
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSpace(line)
		retv = append(retv, line)
	}
	return retv
}
