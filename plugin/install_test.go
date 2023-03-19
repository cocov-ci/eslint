package plugin

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/heyvito/httpie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNodeVersion(t *testing.T) {
	ver := "v12.x"
	root := findParentDir(t)
	fixtures := filepath.Join(root, "plugin", "fixtures")
	pkgJsonPath := filepath.Join(fixtures, pkgJson)
	t.Run("package.json not found", func(t *testing.T) {
		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(fixtures)

		_, err := findNodeVersion(helper.ctx)
		assert.Error(t, err)
		require.EqualError(t, err, errNoPkgJson().Error())
	})

	t.Run("Failed to unmarshall package.json", func(t *testing.T) {
		err := os.WriteFile(pkgJsonPath, []byte{123}, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(fixtures)

		_, err = findNodeVersion(helper.ctx)
		assert.Error(t, err)
	})

	t.Run("Node version not found", func(t *testing.T) {
		err := os.WriteFile(pkgJsonPath, []byte("{}"), os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(fixtures)

		_, err = findNodeVersion(helper.ctx)
		assert.Error(t, err)
		require.EqualError(t, err, errNoVersionFound().Error())
	})

	t.Run("Works as expected with engine field", func(t *testing.T) {
		data := []byte("{\"engine\": {\"node\": \"v12.x\"}}")
		err := os.WriteFile(pkgJsonPath, data, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(fixtures)

		version, err := findNodeVersion(helper.ctx)
		assert.NoError(t, err)
		assert.Equal(t, version, ver)
	})

	t.Run("Works as expected with engines field", func(t *testing.T) {
		data := []byte("{\"engines\": {\"node\": \"v12.x\"}}")
		err := os.WriteFile(pkgJsonPath, data, os.ModePerm)
		require.NoError(t, err)

		t.Cleanup(func() { _ = os.Remove(pkgJsonPath) })

		helper := newTestHelper(t)

		helper.ctx.EXPECT().Workdir().Return(fixtures)

		version, err := findNodeVersion(helper.ctx)
		assert.NoError(t, err)
		assert.Equal(t, version, ver)
	})
}

func TestDetermineNodeVersion(t *testing.T) {
	toEval, versionErr := semver.NewVersion("v9")
	require.NoError(t, versionErr)

	t.Run("Returns a single constraint", func(t *testing.T) {
		v := "v9.x"
		constr, err := determineVersionConstraints(v)
		assert.NoError(t, err)
		assert.Len(t, constr, 1)

		ok := constr.eval(toEval)
		assert.True(t, ok)

		// Semver is able to handle this case as a single constraint
		v = "^8.x || ^10.x"
		constr, err = determineVersionConstraints(v)
		assert.NoError(t, err)
		assert.Len(t, constr, 1)

		ok = constr.eval(toEval)
		assert.False(t, ok)

		ev1, err := semver.NewVersion("v8.9")
		require.NoError(t, err)

		ev2, err := semver.NewVersion("v10.2")
		require.NoError(t, err)

		ok = constr.eval(ev1, ev2)
		assert.True(t, ok)
	})

	t.Run("Returns n constraints", func(t *testing.T) {
		v := ">=v12.x <=v13.4.x "
		constrs, err := determineVersionConstraints(v)
		assert.NoError(t, err)
		assert.Len(t, constrs, 2)

		v = "^8.x  <=10.x"
		constrs, err = determineVersionConstraints(v)
		assert.NoError(t, err)
		assert.Len(t, constrs, 2)

		ok := constrs.eval(toEval)
		assert.True(t, ok)

		ev1, err := semver.NewVersion("v8.9")
		require.NoError(t, err)

		ev2, err := semver.NewVersion("v10.2")
		require.NoError(t, err)

		ok = constrs.eval(ev1, ev2)
		assert.True(t, ok)
	})
}

func TestGetNodeVersionIndex(t *testing.T) {
	t.Run("Fails to retrieve version index", func(t *testing.T) {
		internalErr := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		server := httpie.New(httpie.WithCustom("/", internalErr))
		defer server.Stop()

		helper := newTestHelper(t)

		_, err := getNodeVersionIndex(helper.ctx, server.URL)
		assert.Error(t, err)
	})

	t.Run("Fails to decode response", func(t *testing.T) {
		server := httpie.New(httpie.WithJSON("/", []byte{123}))
		defer server.Stop()

		helper := newTestHelper(t)

		_, err := getNodeVersionIndex(helper.ctx, server.URL)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "json")
	})

	t.Run("Works as expected", func(t *testing.T) {
		helper := newTestHelper(t)
		vi, err := getNodeVersionIndex(helper.ctx, nodeIndex)
		assert.NoError(t, err)
		assert.NotNil(t, vi)
	})
}

func TestFindViableVersion(t *testing.T) {
	vi := &versionIndex{
		versions: []versionInfo{
			{Version: "v12.3"},
			{Version: "v14.2"},
		},
	}

	t.Run("Fails to find a viable version", func(t *testing.T) {
		rawVer := "v11.2"
		constr, err := determineVersionConstraints(rawVer)
		require.NoError(t, err)

		helper := newTestHelper(t)

		_, err = findViableVersion(helper.ctx, rawVer, constr, vi)
		assert.Error(t, err)
	})

	t.Run("Works as expected", func(t *testing.T) {
		rawVer := "v12.3"
		constr, err := determineVersionConstraints(rawVer)
		require.NoError(t, err)

		helper := newTestHelper(t)

		_, err = findViableVersion(helper.ctx, rawVer, constr, vi)
		assert.NoError(t, err)
	})
}
