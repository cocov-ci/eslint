package main

import (
	"github.com/cocov-ci/go-plugin-kit/cocov"

	"github.com/cocov-ci/eslint/plugin"
)

func main() { cocov.Run(plugin.Run) }
