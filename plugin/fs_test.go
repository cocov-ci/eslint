package plugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindRepositories(t *testing.T) {
	root := findRepositoryRoot(t)
	path := filepath.Join(root, "plugin", "fixtures")

	entries, err := os.ReadDir(path)
	assert.NoError(t, err)

	dirCount := 0
	for _, e := range entries {
		if e.IsDir() {
			dirCount++
		}
	}

	repos, err := findRepositories(path)
	require.NoError(t, err)
	assert.NotNil(t, repos)

	ignoredNodeModules := dirCount-len(repos) == 1
	assert.Truef(t, ignoredNodeModules,
		"Should ignore package.json files that are inside node_modules",
	)

}

func TestCheckDependencies(t *testing.T) {
	ver := "v12.x"
	root := findRepositoryRoot(t)
	fixtures := filepath.Join(root, "plugin", "fixtures")
	pkgJsonPath := filepath.Join(fixtures, pkgJson)
	t.Run("Failed to unmarshall package.json", func(t *testing.T) {
		err := os.WriteFile(pkgJsonPath, []byte{123}, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		_, err = checkDependencies(helper.ctx, fixtures)
		assert.Error(t, err)
	})

	t.Run("Node version not found", func(t *testing.T) {
		err := os.WriteFile(pkgJsonPath, []byte("{}"), os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		_, err = checkDependencies(helper.ctx, fixtures)
		assert.Error(t, err)
		require.EqualError(t, err, errNoVersionFound.Error())
	})

	t.Run("Eslint not found as dependency", func(t *testing.T) {
		data := []byte("{\"engines\": {\"node\": \"v12.x\"}}")
		err := os.WriteFile(pkgJsonPath, data, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		_, err = checkDependencies(helper.ctx, fixtures)
		require.EqualError(t, err, errNoEslintDep.Error())
	})

	t.Run("Works as expected with eslint as dependency", func(t *testing.T) {
		data := []byte("{\"engines\": {\"node\": \"v12.x\"}, \"dependencies\": {\"eslint\": \"v8.0\"}}")
		err := os.WriteFile(pkgJsonPath, data, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		version, err := checkDependencies(helper.ctx, fixtures)
		assert.NoError(t, err)
		assert.Equal(t, version, ver)
	})

	t.Run("Works as expected with eslint as development dependency", func(t *testing.T) {
		data := []byte("{\"engines\": {\"node\": \"v12.x\"}, \"devDependencies\": {\"eslint\": \"v8.0\"}}")
		err := os.WriteFile(pkgJsonPath, data, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		version, err := checkDependencies(helper.ctx, fixtures)
		assert.Equal(t, version, ver)
	})
}
