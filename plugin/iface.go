package plugin

import "github.com/cocov-ci/go-plugin-kit/cocov"

type Exec interface {
	Exec(cmd string, args []string, opts *cocov.ExecOpts) ([]byte, error)
	Exec2(string, []string, *cocov.ExecOpts) (stdout, stderr []byte, err error)
}

type ccExec struct{}

func defaultExec() Exec { return ccExec{} }

func (ccExec) Exec2(cmd string, args []string, opts *cocov.ExecOpts) (stdout, stderr []byte, err error) {
	return cocov.Exec2(cmd, args, opts)
}

func (ccExec) Exec(cmd string, args []string, opts *cocov.ExecOpts) ([]byte, error) {
	return cocov.Exec(cmd, args, opts)
}
