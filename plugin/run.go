package plugin

import "github.com/cocov-ci/go-plugin-kit/cocov"

func Run(ctx cocov.Context) error {
	res, err := run(ctx)
	if err != nil {
		return err
	}
	_ = res

	return nil
}

func run(ctx cocov.Context) ([]result, error) {
	exec := defaultExec()
	np, err := installNode(ctx, exec)
	if err != nil {
		return nil, err
	}

	mgr, err := installPkgManager(ctx, exec, np)
	if err != nil {
		return nil, err
	}

	if err = restoreNodeModules(ctx, exec, mgr, ""); err != nil {
		return nil, err
	}

	res, err := runEslint(ctx, exec, mgr, "")
	if err != nil {
		return nil, err
	}

	return res, nil
}
