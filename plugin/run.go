package plugin

import (
	"fmt"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"go.uber.org/zap"
)

func Run(ctx cocov.Context) error {
	out, err := run(ctx)
	if err != nil {
		return err
	}

	sha := ctx.CommitSHA()
	for _, res := range out.Results {
		for _, m := range res.Messages {
			kind, ok := out.kindForRule(m.RuleID)
			if !ok {
				continue
			}

			input := fmt.Sprintf(
				"%s-%d-%s-%s",
				kind.String(), m.Line, res.FilePath, sha,
			)

			id := cocov.SHA1([]byte(input))
			if err = ctx.EmitIssue(kind, res.FilePath, m.Line, m.EndLine, m.Message, id); err != nil {
				ctx.L().Error("Error emitting issue", zap.Error(err))
				return err
			}
		}
	}
	return nil
}

func run(ctx cocov.Context) (*cliOutput, error) {
	exec := defaultExec()
	repos, err := findRepositories(ctx.Workdir())
	if err != nil {
		ctx.L().Error("Failed looking for repositories", zap.Error(err))
		return nil, err
	}

	if len(repos) < 1 {
		ctx.L().Error("Failed to find any package.json files")
		return nil, errNoPkgJson
	}

	out := newCliOutput()
	for _, repo := range repos {
		np, err := installNode(ctx, exec, repo)
		if err != nil {
			return nil, err
		}

		mgr, file, err := installPkgManager(ctx, exec, np, repo)
		if err != nil {
			return nil, err
		}

		if err = restoreNodeModules(ctx, exec, mgr, file, np, repo); err != nil {
			return nil, err
		}

		repoOutput, err := runEslint(ctx, exec, np, repo)
		if err != nil {
			return nil, err
		}
		out.Results = append(out.Results, repoOutput.Results...)

		for k, v := range out.Metadata.RulesMeta {
			out.Metadata.RulesMeta[k] = v
		}

	}

	return out, nil
}
