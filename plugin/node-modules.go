package plugin

import (
	"encoding/json"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

func restoreNodeModules(ctx cocov.Context, e Exec, manager, nodePath string) error {
	envs := map[string]string{"PATH": nodePath}
	opts := &cocov.ExecOpts{Workdir: ctx.Workdir(), Env: envs}
	ctx.L().Info("Restoring node modules")
	stdOut, stdErr, err := e.Exec2(manager, []string{"install"}, opts)
	if err != nil {
		ctx.L().Error("error restoring node modules",
			zap.String("std out", string(stdOut)),
			zap.String("std err", string(stdErr)),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func runEslint(ctx cocov.Context, e Exec, manager, nodePath string) (*cliOutput, error) {
	args := []string{"run", "", "eslint", "-f", "json-with-metadata", "."}

	if manager == yarn {
		args[1] = "-s"
	}

	ctx.L().Info("Running eslint")
	start := time.Now()

	envs := map[string]string{"PATH": nodePath}
	opts := &cocov.ExecOpts{Workdir: ctx.Workdir(), Env: envs}
	stdOut, stdErr, err := e.Exec2(manager, args, opts)
	if err != nil {
		ctx.L().Error("error running eslint: %s",
			zap.String("std out", string(stdOut)),
			zap.String("std err", string(stdErr)),
			zap.Error(err),
		)
		return nil, err
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
