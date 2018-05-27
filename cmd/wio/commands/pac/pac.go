// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio. It used npm as a backend and pushes packages to that
package pac

import (
    "github.com/urfave/cli"
    "path/filepath"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
    "os/exec"
    "bytes"
    "os"
    "strings"
    "wio/cmd/wio/utils/io/log"
    "errors"
    "net/http"
    "regexp"
)

const (
    PUBLISH = "publish"
    GET     = "get"
    UPDATE  = "update"
    COLLECT = "collect"
)

type Pac struct {
    Context *cli.Context
    Type    string
    error
}

// Get context for the command
func (pac Pac) GetContext() (*cli.Context) {
    return pac.Context
}

// Executes the libraries command
func (pac Pac) Execute() {
    directory, err := filepath.Abs(pac.Context.String("dir"))
    commands.RecordError(err, "")

    switch pac.Type {
    case PUBLISH:
        handlePublish(directory)
    case GET:
        handleGet(directory, pac.Context.Bool("clean"))
        break
    case UPDATE:
        handleUpdate(directory)
        break
    case COLLECT:
        handleCollect(directory)
        break
    }
}

// This package.json file will be used to check if dependencies are up to date
func createPrivatePackageJson(directory string, name string, dependencies types.DependenciesTag) {
    log.Norm.Cyan(false, "setting up packages environment ... ")

    // wio path
    wioPath := directory + io.Sep + ".wio"

    // npm config
    npmConfig := types.NpmConfig{}

    // fill all the fields for package.json
    npmConfig.Name = name
    npmConfig.Version = "0.0.1"
    npmConfig.Main = ".wio.js"

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    log.Verb.Verbose(true, "")

    // add dependencies to package.json
    for dependencyName, dependencyValue := range dependencies {
        if !dependencyValue.Vendor {
            npmConfig.Dependencies[dependencyName] = "^" + dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(wioPath+io.Sep+"package.json", &npmConfig); err != nil {
        commands.RecordError(err, "failure")
    }

    log.Norm.Green(true, "success")
}

// publishes the package to npm backend
func handlePublish(directory string) {
    log.Norm.Yellow(true, "Publishing this package using npm server")
    log.Norm.Cyan(false, "verifying project structure ... ")

    // read wio.yml file
    pkgConfig := &types.PkgConfig{}

    io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig)

    if pkgConfig.MainTag.Name == "" {
        commands.RecordError(errors.New("push command only works on project of type \"pkg\""), "failure")
    } else {
        log.Norm.Green(true, "success")
    }

    log.Norm.Cyan(false, "creating package manager files ... ")

    // npm config
    npmConfig := types.NpmConfig{}

    // fill all the fields for package.json
    npmConfig.Name = pkgConfig.MainTag.Name
    npmConfig.Version = pkgConfig.MainTag.Version
    npmConfig.Description = pkgConfig.MainTag.Description
    npmConfig.Repository = pkgConfig.MainTag.Url
    npmConfig.Main = ".wio.js"
    npmConfig.Keywords = utils.AppendIfMissing(pkgConfig.MainTag.Keywords, []string{"c++", "c", "wio"})
    npmConfig.Author = pkgConfig.MainTag.Author
    npmConfig.License = pkgConfig.MainTag.License
    npmConfig.Contributors = pkgConfig.MainTag.Contributors

    npmConfig.Dependencies = make(types.NpmDependencyTag)

    log.Verb.Verbose(true, "")

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if !dependencyValue.Vendor {
            dependencyCheck(directory, dependencyName, dependencyValue.Version)

            npmConfig.Dependencies[dependencyName] = "^" + dependencyValue.Version
        }
    }

    // write package.json file
    if err := io.NormalIO.WriteJson(directory+io.Sep+"package.json", &npmConfig); err != nil {
        commands.RecordError(err, "failure")
    }

    // write .wio.js file so that this is the entry point
    if err := io.NormalIO.WriteFile(directory+io.Sep+".wio.js", []byte("console.log('H"+
        "i!!! Welcome to Wio world!')")); err != nil {
        removeNpmFiles(directory)
        commands.RecordError(err, "failure")
    }

    log.Norm.Green(true, "success")
    log.Norm.Cyan(false, "running publish command ... ")

    // execute cmake command
    npmPublishCommand := exec.Command("npm", "publish")
    npmPublishCommand.Dir = directory

    // Stderr buffer
    cmdErrOutput := &bytes.Buffer{}
    npmPublishCommand.Stderr = cmdErrOutput

    if log.Verb.IsVerbose() {
        npmPublishCommand.Stdout = os.Stdout
    }
    err := npmPublishCommand.Run()
    if err != nil {
        removeNpmFiles(directory)

        // this means user needs to login
        if strings.Contains(cmdErrOutput.String(), "npm adduser") ||
            strings.Contains(cmdErrOutput.String(), "npm login") {
            log.Norm.Yellow(true, "\nNPM backend is used for package management")
            log.Norm.Green(true, "Run \"npm adduser\" to login/create account to publish the package")

            os.Exit(2)
        }

        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    } else {
        log.Norm.Green(true, "success")
    }

    removeNpmFiles(directory)

    log.Verb.Verbose(true, "success")
    log.Norm.Yellow(true, pkgConfig.MainTag.Name+"@"+pkgConfig.MainTag.Version+" published!!")
}

// checks if dependencies are valid wio packages and if they are already pushed
func dependencyCheck(directory string, dependencyName string, dependencyVersion string) {
    log.Verb.Verbose(false, "dependency: checking if "+dependencyName+" package exists ... ")

    resp, err := http.Get("https://www.npmjs.com/package/" + dependencyName + "/v/" + dependencyVersion)
    if err != nil {
        commands.RecordError(err, "failure")
    }

    // dependency does not exist
    if resp.StatusCode == 404 {
        commands.RecordError(errors.New("dependency: \""+dependencyName+"\" package does not exist"),
            "failure")
    }
    resp.Body.Close()

    log.Verb.Verbose(true, "success")
    log.Verb.Verbose(false, "dependency: checking if " + dependencyName + "@" + dependencyVersion+
        " version exists ... ")

    // verify the version by executing npm info command
    npmInfoCommand := exec.Command("npm", "info", dependencyName+"@"+dependencyVersion)
    npmInfoCommand.Dir = directory

    // Stderr buffer
    cmdErrOutput := &bytes.Buffer{}
    cmdOutOutput := &bytes.Buffer{}
    npmInfoCommand.Stderr = cmdErrOutput
    npmInfoCommand.Stdout = cmdOutOutput

    err = npmInfoCommand.Run()
    if err != nil {
        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    }

    // version does not exists
    if cmdOutOutput.String() == "" {
        commands.RecordError(errors.New("dependency: \"" + dependencyName + "@" + dependencyVersion+
            "\" version does not exist"), "failure", )
    } else {
        log.Verb.Verbose(true, "success")

        log.Verb.Verbose(false, "dependency: checking if " + dependencyName + "@" + dependencyVersion+
            " is a valid wio package ... ")

        // check if the package is a wio package by checking C, C++ and wio flags
        pat := regexp.MustCompile(`keywords: .*[\r\n]`)
        s := pat.FindString(cmdOutOutput.String())

        // if wio, c and c++ found, this package is a valid wio package
        if strings.Contains(s, "wio") && strings.Contains(s, "c") && strings.Contains(s, "c++") {
            log.Verb.Verbose(true, "success")
        } else {
            commands.RecordError(errors.New("dependency: \"" + dependencyName + "@" + dependencyVersion+
                "\" is not a wio package"), "failure", )
        }
    }
}

// removes files used for npm publish
func removeNpmFiles(directory string) {
    log.Verb.Verbose(false, "removing npm files ... ")

    if err := os.RemoveAll(directory + io.Sep + "package.json"); err != nil {
        commands.RecordError(err, "failure")
    }

    if err := os.RemoveAll(directory + io.Sep + ".wio.js"); err != nil {
        commands.RecordError(err, "failure")
    }
}

// gets and updates the packages to the versions specified in wio.yml file
func handleGet(directory string, clean bool) {
    log.Norm.Yellow(true, "Getting packages from npm server")

    // clean npm_modules in .wio folder
    if clean {
        log.Norm.Cyan(false, "removing all the pulled packages ... ")

        if err := os.RemoveAll(directory + io.Sep + ".wio" + io.Sep + "node_modules"); err != nil {
            commands.RecordError(err, "failure")
        } else {
            log.Norm.Green(true, "success")
        }
    }

    // read wio.yml file. We gonna use dependencies tag so it does not matter if this is app or pkg
    pkgConfig := &types.PkgConfig{}

    // create private json file
    createPrivatePackageJson(directory, filepath.Base(directory), pkgConfig.DependenciesTag)

    io.NormalIO.ParseYml(directory+io.Sep+"wio.yml", pkgConfig)

    // add dependencies to package.json
    for dependencyName, dependencyValue := range pkgConfig.DependenciesTag {
        if dependencyValue.Vendor {
            continue
        }

        dependencyCheck(directory, dependencyName, dependencyValue.Version)

        log.Norm.Cyan(false, "pulling " + dependencyName + "@" + dependencyValue.Version +
            " package ... ")

        // execute cmake command
        npmInstallCommand := exec.Command("npm", "install", dependencyName + "@" + dependencyValue.Version)
        npmInstallCommand.Dir = directory + io.Sep + ".wio"

        // Stderr buffer
        cmdErrOutput := &bytes.Buffer{}
        npmInstallCommand.Stderr = cmdErrOutput

        if log.Verb.IsVerbose() {
            log.Norm.Verbose(true, "")
            npmInstallCommand.Stdout = os.Stdout
        }
        err := npmInstallCommand.Run()
        if err != nil {
            commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
        }

        log.Norm.Green(true, "success")
    }

    log.Norm.Yellow(true, "All packages pulled successfully")
}

func handleUpdate(directory string) {
    // TODO future versions
    // do npm outdated and see if something needs to be updated
    // then read package.json file and update the version of wio dependency
}

func handleCollect(directory string) {
    // TODO future versions
}
