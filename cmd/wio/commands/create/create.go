// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Creates and initializes a wio project. It also works as an updater when called on already created projects.
package create

import (
    . "wio/cmd/wio/types"
    . "github.com/logrusorgru/aurora"
    "fmt"
    "path/filepath"
)

// Interface for project types and how they are created
type Type interface {
    createStructure() (error)
    printProjectStructure()
    createTemplateProject() (error)
    printNextCommands()
}

// Executes the create command provided configuration packet
func Execute(config ConfigCreate) {
    fmt.Print(Brown("Project Name: "))
    fmt.Println(Cyan(filepath.Base(config.Directory)))
    fmt.Print(Brown("Project Type: "))
    fmt.Println(Cyan(config.AppType))
    fmt.Print(Brown("Project Path: "))
    fmt.Println(Cyan(config.Directory))
    fmt.Println()

    var createType Type = App{config:config}

    if config.AppType == "lib" {
        createType = Lib{config:config}
    }

    fmt.Print(Brown("Creating project structure ..."))
    err := createType.createStructure()

    if err != nil {
        fmt.Println(Red(" [failure]"))
    } else {
        createType.printProjectStructure()
    }

    err = createType.createTemplateProject()

    if err != nil {
        fmt.Println(Red(" [failure]"))
        fmt.Println(Brown("Project creation failed!"))
    } else {
        fmt.Println()
        fmt.Println(Brown("Project has been successfully created and initialized!!"))
        fmt.Println(Brown("Check following commands: "))
        createType.printNextCommands()
    }
}
