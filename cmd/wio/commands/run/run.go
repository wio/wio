// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands package, which contains all the commands provided by the tool.
// Builds, Tests, and Uploads the project to a device
package run
<<<<<<< HEAD
=======

import (
    "github.com/urfave/cli"
    "wio/cmd/wio/commands/build"
)

type Run struct {
    Context *cli.Context
    error
}

// get context for the command
func (run Run) GetContext() (*cli.Context) {
    return run.Context
}

// Runs the build command when cli build option is provided
func (run Run) Execute() {
    // execute build command
    build.RunBuild(run.Context.String("dir"), run.Context.String("target"),
        run.Context.Bool("clean"))

    // execute upload command
}
>>>>>>> More commands and minor fixes (#37)
