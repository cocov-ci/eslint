package plugin

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

func findRepositories(rootPath string) ([]string, error) {
	root := os.DirFS(rootPath)
	var repos []string

	err := fs.WalkDir(root, ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if d.Name() == "package.json" {
				repos = append(repos, filepath.Dir(path))
			}
			return nil
		})

	if err != nil {
		return nil, err
	}

	return repos, nil
}

func checkDependencies(ctx cocov.Context, repoPath string) (string, error) {
	jsonFile := filepath.Join(repoPath, pkgJson)
	f, err := os.ReadFile(jsonFile)
	if err != nil {
		ctx.L().Error("failed to read package.json", zap.Error(err))
		return "", err
	}

	pkg := struct {
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`

		Deps    map[string]string `json:"dependencies"`
		DevDeps map[string]string `json:"devDependencies"`
	}{}

	if err = json.Unmarshal(f, &pkg); err != nil {
		ctx.L().Error("failed to unmarshall package.json", zap.Error(err))
		return "", err
	}

	nodeVersion := pkg.Engines.Node
	if nodeVersion == "" {
		ctx.L().Error(errNoVersionFound.Error())
		return "", errNoVersionFound
	}

	eslintKey := "eslint"
	if _, ok := pkg.Deps[eslintKey]; !ok {
		if _, ok = pkg.DevDeps[eslintKey]; !ok {
			return "", errNoEslintDep
		}
	}

	return nodeVersion, nil
}
