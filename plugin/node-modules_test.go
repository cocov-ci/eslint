package plugin

import (
	"errors"
	"testing"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestoreNodeModules(t *testing.T) {
	wd := "workdir"
	np := "node-path"
	manager := yarn
	opts := &cocov.ExecOpts{Workdir: wd, Env: map[string]string{"PATH": np}}
	t.Run("Fails to restore node modules", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := []byte("something on std out")
		stdErr := []byte("something on std err")
		boom := errors.New("boom")
		helper.exec.EXPECT().
			Exec2(manager, []string{"install"}, opts).
			Return(stdOut, stdErr, boom)

		err := restoreNodeModules(helper.ctx, helper.exec, manager, np)
		require.Error(t, err)
	})

	t.Run("Works as expected", func(t *testing.T) {
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(wd)

		helper.exec.EXPECT().
			Exec2(manager, []string{"install"}, opts).
			Return(nil, nil, nil)

		err := restoreNodeModules(helper.ctx, helper.exec, manager, np)
		require.NoError(t, err)
	})
}

func TestRunEslint(t *testing.T) {
	wd := "workdir"
	np := "node-path"
	manager := yarn
	opts := &cocov.ExecOpts{Workdir: wd, Env: map[string]string{"PATH": np}}
	t.Run("Fails running eslint", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := []byte("something on std out")
		stdErr := []byte("something on std err")
		boom := errors.New("boom")

		helper.exec.EXPECT().
			Exec2(manager, []string{"run", "-s", "eslint", "-f", "json-with-metadata", "."}, opts).
			Return(stdOut, stdErr, boom)

		_, err := runEslint(helper.ctx, helper.exec, manager, np)
		require.Error(t, err)
	})

	t.Run("Fails unmarshalling output", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := []byte("123")
		stdErr := []byte("something went wrong")

		helper.exec.EXPECT().
			Exec2(manager, []string{"run", "-s", "eslint", "-f", "json-with-metadata", "."}, opts).
			Return(stdOut, stdErr, nil)

		_, err := runEslint(helper.ctx, helper.exec, manager, np)
		require.Error(t, err)
		assert.ErrorContains(t, err, "json")
	})

	t.Run("Works as expected", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := validOutput(t)

		helper.exec.EXPECT().
			Exec2(manager, []string{"run", "-s", "eslint", "-f", "json-with-metadata", "."}, opts).
			Return(stdOut, nil, nil)

		out, err := runEslint(helper.ctx, helper.exec, manager, np)
		require.NoError(t, err)
		assert.NotNil(t, out)
	})
}

func validOutput(t *testing.T) []byte {
	root := findRepositoryRoot(t)
	dataPath := filepath.Join(root, "plugin", "fixtures", "out.json")

	data, err := os.ReadFile(dataPath)
	require.NoError(t, err)
	return data
}
