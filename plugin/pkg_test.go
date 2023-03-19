package plugin

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallPkgManager(t *testing.T) {
	root := findParentDir(t)
	fixtureRoot := filepath.Join(root, "plugin", "fixtures")

	t.Run("Lockfile not found", func(t *testing.T) {
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(fixtureRoot)

		_, err := findLockFile(helper.ctx)
		assert.Error(t, err)
	})

	t.Run("Founds npm", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "npm")
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(p)

		mgr, err := findLockFile(helper.ctx)
		assert.NoError(t, err)
		assert.Equal(t, mgr, npm)
	})

	t.Run("Founds yarn", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "yarn")
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(p)

		mgr, err := findLockFile(helper.ctx)
		assert.NoError(t, err)
		assert.Equal(t, mgr, yarn)
	})

	t.Run("Founds pnpm", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "pnpm")
		helper := newTestHelper(t)
		helper.ctx.EXPECT().Workdir().Return(p)

		mgr, err := findLockFile(helper.ctx)
		assert.NoError(t, err)
		assert.Equal(t, mgr, pnpm)
	})
}
