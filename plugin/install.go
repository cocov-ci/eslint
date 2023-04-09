package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	version, err := findNodeVersion(ctx)
	if err != nil {
		return "", err
	}

	cachedVersion := toolCacheVersion(version)
	if ok := ctx.LoadToolCache(cachedVersion, nodePath); ok {
		return nodePath, nil
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

	binPath, err := untar(ctx, exec, zip)
	if err != nil {
		return "", err
	}

	ctx.StoreToolCache(cachedVersion, binPath)
	rawPath := os.Getenv("PATH")
	np := fmt.Sprintf("%s:%s", binPath, rawPath)

	return np, nil
}

func findNodeVersion(ctx cocov.Context) (string, error) {
	pkgPath := filepath.Join(ctx.Workdir(), pkgJson)
	f, err := os.ReadFile(pkgPath)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.L().Error("package.json not found", zap.Error(err))
			return "", errNoPkgJson
		}

		ctx.L().Error("failed to read package.json", zap.Error(err))
		return "", err
	}

	pkg := struct {
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`
	}{}

	if err = json.Unmarshal(f, &pkg); err != nil {
		ctx.L().Error("failed to unmarshall package.json", zap.Error(err))
		return "", err
	}

	enginesVersion := pkg.Engines.Node
	if enginesVersion == "" {
		ctx.L().Error(errNoVersionFound.Error())
		return "", errNoVersionFound
	}

	return enginesVersion, nil
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

func untar(ctx cocov.Context, e Exec, filePath string) (string, error) {
	args := []string{"zxf", filePath, "--strip", "1", "-C", nodePath}
	if _, err := e.Exec("tar", args, nil); err != nil {
		ctx.L().Error("error extracting downloaded file", zap.Error(err))
		return "", err
	}

	_ = os.Remove(filePath)

	return filepath.Join(filePath, "bin"), nil
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

func toolCacheVersion(version string) string {
	return fmt.Sprintf("node-%s-linux-x64", version)
}
