package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
	"os/exec"
	"path/filepath"
	"time"
)

func runEslint(ctx cocov.Context, e Exec, nodePath string) (*cliOutput, error) {
	wd := ctx.Workdir()
	eslintPath := filepath.Join(wd, "node_modules", ".bin", "eslint")
	args := []string{"-f", "json-with-metadata", "."}

	ctx.L().Info("Running eslint")
	start := time.Now()

	envs := map[string]string{"PATH": nodePath}
	opts := &cocov.ExecOpts{Env: envs}
	stdOut, stdErr, err := e.Exec2(eslintPath, args, opts)
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			if execErr.ExitCode() != 1 {
				ctx.L().Error("eslint exited with unexpected status",
					zap.Int("status", execErr.ExitCode()),
					zap.String("std err", string(stdErr)),
					zap.Error(err),
				)
				return nil, err
			}

		} else if !ok {
			ctx.L().Error("error running eslint",
				zap.String("std err", string(stdErr)),
				zap.Error(err),
			)
			return nil, err
		}
	}

	msg := fmt.Sprintf("Running eslint took %s seconds", time.Since(start))
	ctx.L().Info(msg)

	out := &cliOutput{}
	if err = json.Unmarshal(stdOut, out); err != nil {
		ctx.L().Error("failed to unmarshall output",
			zap.Error(err),
		)
		return nil, err
	}
	return out, nil
}
