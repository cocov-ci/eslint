package plugin

import "github.com/cocov-ci/go-plugin-kit/cocov"

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
