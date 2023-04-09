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

			kind, ok := rules[m.RuleID]
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
	np, err := installNode(ctx, exec)
	if err != nil {
		return nil, err
	}

	mgr, file, err := installPkgManager(ctx, exec, np)
	if err != nil {
		return nil, err
	}

	if err = restoreNodeModules(ctx, exec, mgr, file, np); err != nil {
		return nil, err
	}

	res, err := runEslint(ctx, exec, np)
	if err != nil {
		return nil, err
	}

	return res, nil
}
