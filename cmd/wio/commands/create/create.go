// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Creates and initializes a wio project. It also works as an updater when called on already created projects.
package create

import (
    . "wio/cmd/wio/types"
    . "wio/cmd/wio/io"
    "path/filepath"
    "os"
)

// Interface for project types and how they are created
type Type interface {
    createStructure() (error)
    printProjectStructure()
    createTemplateProject() (error)
    printNextCommands()
}

// Executes the create command provided configuration packet
func Execute(config ConfigCreate, ioData IOData) {
    Norm.Yellow("Project Name: ")
    Norm.Cyan(filepath.Base(config.Directory) + "\n")
    Norm.Yellow("Project Type: ")
    Norm.Cyan(config.AppType + "\n")
    Norm.Yellow("Project Path: ")
    Norm.Cyan(config.Directory + "\n")
    Norm.White("\n")

    var createType Type = App{config:config, ioData:ioData}

    if config.AppType == "lib" {
        createType = Lib{config:config, ioData:ioData}
    }

    Norm.Yellow("Creating project structure ... ")
    err := createType.createStructure()

    if err != nil {
        Norm.Red( "[failure]\n")
        Verb.Error(err.Error() + "\n")
        os.Exit(2)
    } else {
        Norm.Green("[success]\n")
        createType.printProjectStructure()
    }

    Norm.White("\n")
    Norm.Yellow("Creating template project ... ")
    err = createType.createTemplateProject()

    if err != nil {
        Norm.Red("[failure]\n")
        Verb.Error(err.Error() + "\n")
    } else {
        Norm.Green("[success]\n")
        Norm.White("\n")
        Norm.Yellow( "Project has been successfully created and initialized!!\n")
        Norm.Green("Check following commands: \n")
        createType.printNextCommands()
    }
}
