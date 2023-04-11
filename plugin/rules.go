package plugin

import (
	"strings"

	"github.com/cocov-ci/go-plugin-kit/cocov"
)

//go:generate go run ../generator/genrules.go

type message struct {
	RuleID  string `json:"ruleId"`
	Message string `json:"message"`
	Line    uint   `json:"line"`
	EndLine uint   `json:"endLine"`
}

type result struct {
	FilePath string    `json:"filePath"`
	Messages []message `json:"messages"`
}

type metadata struct {
	RulesMeta map[string]metadataInfo
}

type metadataInfo struct {
	Type string `json:"type"`
}

type cliOutput struct {
	Results  []result `json:"results"`
	Metadata metadata `json:"metadata"`
}

func newCliOutput() *cliOutput { return &cliOutput{} }

func (c *cliOutput) kindForRule(rule string) (cocov.IssueKind, bool) {
	v, ok := c.Metadata.RulesMeta[rule]
	if ok {
		switch v.Type {
		case "problem":
			return cocov.IssueKindBug, true
		case "suggestion":
			return cocov.IssueKindConvention, true
		case "layout":
			return cocov.IssueKindStyle, true
		}
	}

	if strings.Contains(rule, "/") {
		split := strings.Split(rule, "/")
		rule = split[len(split)-1]
	}

	kind, ok := rules[rule]
	if ok {
		return kind, true
	}

	return 0, false
}
