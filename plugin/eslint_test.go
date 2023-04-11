package plugin

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/cocov-ci/go-plugin-kit/cocov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunEslint(t *testing.T) {
	wd := "workdir"
	np := "node-path"
	eslintPath := filepath.Join(wd, "node_modules", ".bin", "eslint")
	args := []string{"-f", "json-with-metadata", "--quiet", "."}
	opts := &cocov.ExecOpts{Env: map[string]string{"PATH": np}}

	t.Run("Fails running eslint", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd).AnyTimes()

		stdOut := []byte("something on std out")
		stdErr := []byte("something on std err")
		boom := errors.New("boom")

		helper.exec.EXPECT().
			Exec2(eslintPath, args, opts).
			Return(stdOut, stdErr, boom)

		_, err := runEslint(helper.ctx, helper.exec, np)
		require.Error(t, err)
	})

	t.Run("Fails unmarshalling output", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := []byte("123")
		stdErr := []byte("something went wrong")

		helper.exec.EXPECT().
			Exec2(eslintPath, args, opts).
			Return(stdOut, stdErr, nil)

		_, err := runEslint(helper.ctx, helper.exec, np)
		require.Error(t, err)
		assert.ErrorContains(t, err, "json")
	})

	t.Run("Works as expected", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(wd)

		stdOut := validOutput(t)

		helper.exec.EXPECT().
			Exec2(eslintPath, args, opts).
			Return(stdOut, nil, nil)

		out, err := runEslint(helper.ctx, helper.exec, np)
		require.NoError(t, err)
		assert.NotNil(t, out)
	})
}
