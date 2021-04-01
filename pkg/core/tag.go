package core

import (
	"bytes"
	"html/template"
	"time"

	"github.com/spiarh/gojo/pkg/util"
)

func BuildTag(facts []*Fact, tagFormat, imageDir string) (string, error) {
	factsMap := make(map[string]string)
	for _, f := range facts {
		factsMap[f.Name] = f.Value
	}

	// Date
	now := time.Now().Format("20060102150405")
	factsMap["date"] = now

	// Git
	gitCommit, err := util.GetGitHeadHash(imageDir)
	if err != nil {
		return "", err
	}
	factsMap["gitCommit"] = gitCommit[:8]

	tmpl, err := template.New("Tag").Option("missingkey=error").Parse(tagFormat)
	if err != nil {
		return "", err
	}

	var tag bytes.Buffer
	if err := tmpl.Execute(&tag, factsMap); err != nil {
		return "", err
	}

	return tag.String(), nil
}
