package plugin

import (
	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

func restoreNodeModules(ctx cocov.Context, e Exec, manager, file, nodePath string) error {
	nodeModules := "node_modules"
	artifactKeys := []string{pkgJson, file}

	if _, err := ctx.LoadArtifactCache(artifactKeys, nodeModules); err != nil {
		ctx.L().Error("Error restoring cache artifacts", zap.Error(err))
		return err
	}

	envs := map[string]string{"PATH": nodePath}
	opts := &cocov.ExecOpts{Workdir: ctx.Workdir(), Env: envs}

	ctx.L().Info("Restoring node modules", zap.String("package manager", manager))
	stdOut, stdErr, err := e.Exec2(manager, []string{"install"}, opts)
	if err != nil {
		ctx.L().Error("error restoring node modules",
			zap.String("std out", string(stdOut)),
			zap.String("std err", string(stdErr)),
			zap.Error(err),
		)
		return err
	}

	if err = ctx.StoreArtifactCache(artifactKeys, nodeModules); err != nil {
		ctx.L().Error("Error storing cache artifact", zap.Error(err))
		return err
	}

	return nil
}
