package plugin

import (
	"fmt"
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
	"package-lock.json": npm,
	"pnpm-lock.yaml":    pnpm,
	"yarn.lock":         yarn,
}

func installPkgManager(ctx cocov.Context, e Exec) (string, error) {
	mgr, err := findLockFile(ctx)
	if err != nil {
		return "", err
	}

	ctx.L().Info(fmt.Sprintf("Using %s as package manager", mgr))

	_, err = e.Exec(npm, []string{"install", "-g", mgr}, nil)
	if err != nil {
		ctx.L().Error("failed to install manager", zap.Error(err))
		return "", err
	}

	return mgr, nil

}

func findLockFile(ctx cocov.Context) (string, error) {
	entries, err := os.ReadDir(ctx.Workdir())
	if err != nil {
		ctx.L().Error("error looking for lockfile",
			zap.String("path", ctx.Workdir()),
			zap.Error(err),
		)
		return "", err
	}

	for _, e := range entries {
		if !e.IsDir() {
			v, ok := managers[e.Name()]
			if ok {
				return v, nil
			}
		}
	}

	return "", errLockFileNotFound()
}
