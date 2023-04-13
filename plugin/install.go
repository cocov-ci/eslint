package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/cocov-ci/go-plugin-kit/cocov"
	"github.com/levigross/grequests"
	"go.uber.org/zap"
)

const (
	nodeIndex = "https://nodejs.org/dist/index.json"
	nodePath  = "/cocov/node"
	pkgJson   = "package.json"
)

type versionInfo struct {
	Version string   `json:"version"`
	Npm     string   `json:"npm"`
	Files   []string `json:"files"`
}

type versionIndex struct {
	versions []versionInfo
}

func installNode(ctx cocov.Context, exec Exec) (string, error) {
	rawPath := os.Getenv("PATH")
	binPath := path.Join(nodePath, "bin")
	np := fmt.Sprintf("%s:%s", binPath, rawPath)

	version, err := checkDependencies(ctx)
	if err != nil {
		return "", err
	}

	tck := toolCacheKey(version)
	if ok := ctx.LoadToolCache(tck, nodePath); ok {
		return np, nil
	}

	consts, err := determineVersionConstraints(version)
	if err != nil {
		return "", err
	}

	index, err := getNodeVersionIndex(ctx, nodeIndex)
	if err != nil {
		return "", err
	}

	availableVersion, err := findViableVersion(ctx, version, consts, index)
	if err != nil {
		return "", err
	}

	url := downloadURL(availableVersion)
	zip, err := downloadNode(ctx, url)
	if err != nil {
		return "", err
	}

	err = untar(ctx, exec, zip)
	if err != nil {
		return "", err
	}

	ctx.StoreToolCache(tck, nodePath)

	return np, nil
}

func determineVersionConstraints(version string) (constraints, error) {
	v, err := semver.NewConstraint(version)
	if err == nil {
		return constraints{v}, nil
	}

	p := newParser()

	return p.parse(version)
}

func getNodeVersionIndex(ctx cocov.Context, url string) (*versionIndex, error) {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		ctx.L().Error("failed to retrieve node version index", zap.Error(err))
		return nil, err
	}

	ji := versionIndex{versions: []versionInfo{}}
	if err = resp.JSON(&ji.versions); err != nil {
		ctx.L().Error("failed to decode response", zap.Error(err))
		return nil, err
	}

	return &ji, nil
}

func findViableVersion(ctx cocov.Context, base string, c constraints, index *versionIndex) (*semver.Version, error) {
	for _, rawInfo := range index.versions {
		v, err := semver.NewVersion(rawInfo.Version)
		if err != nil {
			ctx.L().Error("failed to build semver version using node index",
				zap.String("version used", rawInfo.Version),
				zap.Error(err),
			)
			return nil, err
		}

		if ok := c.eval(v); ok {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no compatible versions found for %s", base)
}

func downloadNode(ctx cocov.Context, url string) (string, error) {
	fileName := "node.tar.gz"

	if err := os.Mkdir(nodePath, os.ModePerm); err != nil {
		ctx.L().Error("error creating directory",
			zap.String("path", nodePath),
			zap.Error(err),
		)
		return "", err
	}

	ctx.L().Info("downloading node", zap.String("url", url))
	resp, err := grequests.Get(url, nil)
	if err != nil {
		ctx.L().Error("error downloading node", zap.Error(err))
		return "", err
	}

	defer resp.Close()

	tarPath := filepath.Join(nodePath, fileName)
	if err = resp.DownloadToFile(tarPath); err != nil {
		ctx.L().Error("error writing downloaded node to file", zap.Error(err))
		return "", err
	}

	return tarPath, nil
}

func untar(ctx cocov.Context, e Exec, filePath string) error {
	args := []string{"zxf", filePath, "--strip", "1", "-C", nodePath}
	if _, err := e.Exec("tar", args, nil); err != nil {
		ctx.L().Error("error extracting downloaded file", zap.Error(err))
		return err
	}

	_ = os.Remove(filePath)

	return nil
}

func errLockFileNotFound() error {
	mgrs := make([]string, 0, len(managers))
	for k := range managers {
		mgrs = append(mgrs, k)
	}

	msg := fmt.Sprintf("lock file not found. supported are: %s",
		strings.Join(mgrs, ", "))

	msg = strings.TrimSpace(msg)
	msg = msg + "."

	return fmt.Errorf(msg)
}

func downloadURL(version *semver.Version) string {
	strVersion := version.String()
	if !strings.HasPrefix(strVersion, "v") {
		strVersion = "v" + strVersion
	}

	return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-linux-x64.tar.gz",
		strVersion, strVersion)
}

var errNoPkgJson = errors.New("package.json not found")
var errNoVersionFound = errors.New("failed to determine node version using package.json")
var errNoEslintDep = errors.New("eslint not found as a project dependency")

func toolCacheKey(version string) string {
	return fmt.Sprintf("node-%s-linux-x64", version)
}
