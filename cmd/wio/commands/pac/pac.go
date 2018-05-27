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
        handleGet(directory)
        break
    case UPDATE:
        handleUpdate(directory)
        break
    case COLLECT:
        handleCollect(directory)
        break
    }
}

// publishes the package to npm backend
func handlePublish(directory string) {
    log.Norm.Yellow(true, "Publishing this package (we use npm backend)")
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

        log.Verb.Verbose(false, "checking if "+dependencyName+" dependency exists ... ")

        resp, err := http.Get("https://www.npmjs.com/package/" + dependencyName + "/v/" + dependencyValue.Version)
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
        log.Verb.Verbose(false, "checking if " + dependencyName + "@" + dependencyValue.Version+
            " version exists ... ")

        // verify the version by executing npm info command
        npmInfoCommand := exec.Command("npm", "info", dependencyName+"@"+dependencyValue.Version)
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
            commands.RecordError(errors.New("dependency: \"" + dependencyName + "@" + dependencyValue.Version+
                "\" version does not exist"), "failure", )
        } else {
            log.Verb.Verbose(true, "success")
        }

        npmConfig.Dependencies[dependencyName] = "^" + dependencyValue.Version
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

func handleGet(directory string) {
}

func handleUpdate(directory string) {
}

func handleCollect(directory string) {
}
