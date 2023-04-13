package plugin

import (
	"path/filepath"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

func restoreNodeModules(ctx cocov.Context, e Exec, manager, file, nodePath, repoPath string) error {
	file = filepath.Join(repoPath, file)
	repoJsonFile := filepath.Join(repoPath, pkgJson)
	nodeModules := filepath.Join(repoPath, "node_modules")
	artifactKeys := []string{repoJsonFile, file}

	if _, err := ctx.LoadArtifactCache(artifactKeys, nodeModules); err != nil {
		ctx.L().Error("Error restoring cache artifacts", zap.Error(err))
		return err
	}

	envs := map[string]string{"PATH": nodePath}
	opts := &cocov.ExecOpts{Workdir: repoPath, Env: envs}

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
