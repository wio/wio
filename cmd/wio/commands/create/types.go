// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Types for all the things being used in this package

package create

import . "wio/cmd/wio/utils/types"

// Interface for Types of Applications
type ProjectTypes interface {
    createStructure() (error)
    printProjectStructure()
    createTemplateProject() (error)
    printNextCommands()
}

// All the data for the project (app)
type App struct {
    args *CliArgs
}


// All the data for the project (lib)
type Lib struct {
    args *CliArgs
}

type Data struct {
    Id string
    Src string
    Des string
    Override bool
}

type Paths struct {
    Paths []Data
}