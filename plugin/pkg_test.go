package plugin

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallPkgManager(t *testing.T) {
	root := findRepositoryRoot(t)
	fixtureRoot := filepath.Join(root, "plugin", "fixtures")

	t.Run("Lockfile not found", func(t *testing.T) {
		helper := newTestHelper(t)
		_, _, err := findLockFile(helper.ctx, fixtureRoot)
		assert.Error(t, err)
	})

	t.Run("Founds npm", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "npm")
		helper := newTestHelper(t)
		mgr, _, err := findLockFile(helper.ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, mgr, npm)
	})

	t.Run("Founds yarn", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "yarn")
		helper := newTestHelper(t)
		mgr, _, err := findLockFile(helper.ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, mgr, yarn)
	})

	t.Run("Founds pnpm", func(t *testing.T) {
		p := filepath.Join(fixtureRoot, "pnpm")
		helper := newTestHelper(t)
		mgr, _, err := findLockFile(helper.ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, mgr, pnpm)
	})
}
