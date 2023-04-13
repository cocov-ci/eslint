package plugin

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepositories(t *testing.T) {
	root := findRepositoryRoot(t)
	path := filepath.Join(root, "plugin", "fixtures")

	repos, err := findRepositories(path)
	require.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 3)
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