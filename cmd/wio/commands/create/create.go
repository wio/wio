// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Creates and initializes a wio project. It also works as an updater when called on already created projects.
package create

import (
    . "wio/cmd/wio/utils/types"
    . "wio/cmd/wio/utils/io"
    "path/filepath"
    "os"
)

// Executes the create command provided configuration packet
func Execute(args CliArgs) {
    Norm.Yellow("Project Name: ")
    Norm.Cyan(filepath.Base(args.Directory) + "\n")
    Norm.Yellow("Project Type: ")
    Norm.Cyan(args.AppType + "\n")
    Norm.Yellow("Project Path: ")
    Norm.Cyan(args.Directory + "\n")
    Norm.White("\n")

    var projectType ProjectTypes = App{args:&args}

    if args.AppType == "lib" {
        projectType = Lib{args:&args}
    }

    Norm.Yellow("Creating project structure ... ")
    err := projectType.createStructure()

    if err != nil {
        Norm.Red( "[failure]\n")
        Verb.Error(err.Error() + "\n")
        os.Exit(2)
    } else {
        Norm.Green("[success]\n")
        projectType.printProjectStructure()
    }

    Norm.White("\n")
    Norm.Yellow("Creating template project ... ")
    err = projectType.createTemplateProject()

    if err != nil {
        Norm.Red("[failure]\n")
        Verb.Error(err.Error() + "\n")
    } else {
        Norm.Green("[success]\n")
        Norm.White("\n")
        Norm.Yellow( "Project has been successfully created and initialized!!\n")
        Norm.Green("Check following commands: \n")
        projectType.printNextCommands()
    }
}
