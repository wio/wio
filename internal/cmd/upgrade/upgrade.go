package upgrade

import (
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"wio/internal/config/meta"
	"wio/internal/config/root"
	"wio/pkg/log"
	"wio/pkg/npm/semver"
	"wio/pkg/util"
	"wio/pkg/util/sys"
	"wio/pkg/util/template"

	"github.com/dhillondeep/go-getrelease"
	"github.com/urfave/cli"
)

type Upgrade struct {
	Context *cli.Context
}

// get context for the command
func (upgrade Upgrade) GetContext() *cli.Context {
	return upgrade.Context
}

// Runs the build command when cli build option is provided
func (upgrade Upgrade) Execute() error {
	forceFlag := upgrade.Context.Bool("force")
	githubClient := getrelease.NewGithubClient(nil)

	preSetup := func() (string, error) {
		cachePath := sys.Path(root.GetUpdatePath(), sys.Cache)

		if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
			return "", util.Error("error creating cache folder")
		}
		return cachePath, nil
	}

	releaseAssetName := template.Replace(wioAssetName, map[string]string{
		"platform": currOsMapping[runtime.GOOS],
		"arch":     currArchMapping[runtime.GOARCH],
		"format":   formatMapping[runtime.GOOS],
	})

	options := []getrelease.Options{
		getrelease.WithArchive("false"), getrelease.WithChecksum(checksumFile),
		getrelease.WithProgress(defaultProgressBar),
	}

	var newVersion string
	var assetName string
	var assetPath string

	execName := template.Replace(wioExecName, map[string]string{
		"extension": extensionMapping[runtime.GOOS],
	})

	if upgrade.Context.NArg() > 0 {
		checksumFileToUse := checksumFile

		version := upgrade.Context.Args()[0]

		versionToUpgradeSem := semver.Parse(version)
		if versionToUpgradeSem == nil {
			return util.Error("wio version %s is invalid", version)
		}

		if !forceFlag && versionToUpgradeSem.LT(*semver.Parse(compatibilityLowerBound)) {
			return util.Error(
				"for compatibility, wio can only be upgraded to versions >= %s. Use --force flag to override",
				compatibilityLowerBound)
		}

		if versionToUpgradeSem.LT(*semver.Parse(assetNameChangeVersion)) {
			checksumFileToUse = preChecksumFile
			releaseAssetName = template.Replace(preWioAssetName, map[string]string{
				"platform": preOsMapping[runtime.GOOS],
				"arch":     preArchMapping[runtime.GOARCH],
				"format":   formatMapping[runtime.GOOS],
			})
			execName = template.Replace(preWioExecName, map[string]string{
				"platform":  preOsMapping[runtime.GOOS],
				"arch":      preArchMapping[runtime.GOARCH],
				"extension": extensionMapping[runtime.GOOS],
			})
		}

		tagName := fmt.Sprintf("v%s", versionToUpgradeSem.String())

		log.Info(log.Cyan, "Downloading the release for version: ")
		log.Info(log.Green, versionToUpgradeSem.String())
		log.Infoln(log.Cyan, "")

		path, err := preSetup()
		if err != nil {
			return err
		}

		// version specified so use that
		assetName, err = getrelease.GetTagAsset(githubClient, path, releaseAssetName,
			"wio", "wio", tagName, append(options, getrelease.WithChecksum(checksumFileToUse))...)
		if err != nil {
			return err
		}

		assetPath = sys.Path(path, assetName)
		newVersion = versionToUpgradeSem.String()
	} else {
		log.Infoln(log.Cyan, "Downloading the latest version")

		path, err := preSetup()
		if err != nil {
			return err
		}

		// no version specified so we gotta pull the latest version
		assetName, err = getrelease.GetLatestAsset(githubClient, path, releaseAssetName,
			"wio", "wio", options...)
		if err != nil {
			return err
		}

		assetPath = sys.Path(path, assetName)

		re, err := regexp.Compile(`[0-9].[0-9].[0.9]`)
		if err != nil {
			return err
		}
		newVersion = re.FindStringSubmatch(assetName)[0]
	}

	log.Info(log.Cyan, "Extracting release file: ")
	log.Infoln(log.Green, assetName)

	extractionPath, err := ioutil.TempDir(os.TempDir(), assetName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(extractionPath)

	if err := archiver.Unarchive(assetPath, extractionPath); err != nil {
		return err
	}

	log.Info(log.Cyan, "Updating wio: ")
	log.Info(log.Green, meta.Version)
	log.Info(log.Cyan, " => ")
	log.Infoln(log.Green, newVersion)

	newWioExec, err := os.Open(sys.Path(extractionPath, execName))
	if err != nil {
		return util.Error("release asset %s not downloaded or corrupted", assetName)
	}

	if err = update.Apply(newWioExec, update.Options{}); err != nil {
		return util.Error("failed to upgrade wio version to %s", newVersion)
	} else {
		wioRootConfig := &root.WioRootConfig{}
		if err := sys.NormalIO.ParseJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
			return util.Error("wio upgraded to %s but, an error occurred and it's not complete", newVersion)
		}
		wioRootConfig.Updated = true
		if err := sys.NormalIO.WriteJson(root.GetConfigFilePath(), wioRootConfig); err != nil {
			return util.Error("wio upgraded to %s but, an error occurred and it's not complete", newVersion)
		}
	}

	log.Infoln(log.Green, "|>Done")

	return nil
}
