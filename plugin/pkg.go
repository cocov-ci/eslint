package plugin

import (
	"os"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

const (
	npm  = "npm"
	pnpm = "pnpm"
	yarn = "yarn"
)

var managers = map[string]string{
	"pnpm-lock.yaml":    pnpm,
	"yarn.lock":         yarn,
	"package-lock.json": npm,
}

func installPkgManager(ctx cocov.Context, e Exec, nodePath string) (string, string, error) {
	mgr, file, err := findLockFile(ctx)
	if err != nil {
		return "", "", err
	}

	if mgr == npm {
		return npm, file, nil
	}

	opts := &cocov.ExecOpts{Env: map[string]string{"PATH": nodePath}}
	_, err = e.Exec(npm, []string{"install", "-g", mgr}, opts)
	if err != nil {
		ctx.L().Error("failed to install manager", zap.Error(err))
		return "", "", err
	}

	return mgr, file, nil

}

func findLockFile(ctx cocov.Context) (string, string, error) {
	entries, err := os.ReadDir(ctx.Workdir())
	if err != nil {
		ctx.L().Error("error looking for lockfile",
			zap.String("path", ctx.Workdir()),
			zap.Error(err),
		)
		return "", "", err
	}

	for _, e := range entries {
		if !e.IsDir() {
			v, ok := managers[e.Name()]
			if ok {
				return v, e.Name(), nil
			}
		}
	}

	return "", "", errLockFileNotFound()
}
