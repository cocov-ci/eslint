package plugin

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
