// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands package, which contains all the commands provided by the tool.
// Cleans all the build files for the project
package clean

import (
    "github.com/urfave/cli"
    "path/filepath"
    "wio/cmd/wio/commands"
    "os"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
)

type Clean struct {
    Context *cli.Context
    error
}

// get context for the command
func (clean Clean) GetContext() (*cli.Context) {
    return clean.Context
}

// Runs the build command when cli build option is provided
func (clean Clean) Execute() {
    RunClean(clean.Context.String("dir"), clean.Context.String("target"),
        clean.Context.IsSet("target"))
}

func RunClean(directoryCli string, targetCli string, targetProvided bool) {
    directory, err := filepath.Abs(directoryCli)
    commands.RecordError(err, "")

    targetToCleanName := targetCli

    if !targetProvided {
        // clean everything
        commands.RecordError(os.RemoveAll(directory+io.Sep+".wio"+io.Sep+"build"+io.Sep+"targets"), "")
    } else {
        // clean an individual target
        targetPath := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets" + io.Sep + targetToCleanName

        if utils.PathExists(targetPath) {
            // delete the target
            commands.RecordError(os.RemoveAll(targetPath), "")
        }
    }
}
