package plugin

import "github.com/cocov-ci/go-plugin-kit/cocov"

type message struct {
	RuleID  string `json:"ruleId"`
	Message string `json:"message"`
	Line    uint16 `json:"line"`
	EndLine uint16 `json:"endLine"`
}

type result struct {
	FilePath string    `json:"filePath"`
	Messages []message `json:"messages"`
}

// TODO: classify Eslint rules
var rules = map[string]cocov.IssueKind{}
