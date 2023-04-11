package plugin

import (
	"errors"
	"testing"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"github.com/stretchr/testify/require"
)

func TestRestoreNodeModules(t *testing.T) {
	wd := "workdir"
	np := "node-path"
	nodeModules := "node_modules"
	lockFile := "yarn.lock"

	artifactKeys := []string{"package.json", lockFile}
	manager := yarn
	opts := &cocov.ExecOpts{Workdir: wd, Env: map[string]string{"PATH": np}}
	t.Run("Fails to restore node modules", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		helper.ctx.EXPECT().LoadArtifactCache(artifactKeys, nodeModules)

		stdOut := []byte("something on std out")
		stdErr := []byte("something on std err")
		boom := errors.New("boom")
		helper.exec.EXPECT().
			Exec2(manager, []string{"install"}, opts).
			Return(stdOut, stdErr, boom)

		err := restoreNodeModules(helper.ctx, helper.exec, manager, lockFile, np)
		require.Error(t, err)
	})

	t.Run("Works as expected", func(t *testing.T) {
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(wd)

		helper.ctx.EXPECT().LoadArtifactCache(artifactKeys, nodeModules)

		helper.exec.EXPECT().
			Exec2(manager, []string{"install"}, opts).
			Return(nil, nil, nil)

		helper.ctx.EXPECT().StoreArtifactCache(artifactKeys, nodeModules)

		err := restoreNodeModules(helper.ctx, helper.exec, manager, lockFile, np)
		require.NoError(t, err)
	})
}
