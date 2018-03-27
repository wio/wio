// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates an executable application
package create

import (
    "os"
    "fmt"
    "path/filepath"

    . "wio/cmd/wio/types"
    . "github.com/logrusorgru/aurora"
)

type App struct {
    config ConfigCreate
}

// Creates project structure for application type
func (app App) createStructure() (error) {
    srcPath := app.config.Directory + string(filepath.Separator) + "src"
    libPath := app.config.Directory + string(filepath.Separator) + "lib"

    err := os.MkdirAll(srcPath, os.ModePerm)
    if err != nil {
        return err
    }

    err = os.MkdirAll(libPath, os.ModePerm)
    if err != nil {
        return err
    }

    if app.config.Tests {
        testPath := app.config.Directory + string(filepath.Separator) + "test"
        err = os.MkdirAll(testPath, os.ModePerm)
        if err != nil {
            return err
        }
    }

    return nil
}

// Prints the project structure for application type
func (app App) printProjectStructure() {
    fmt.Println()
    fmt.Println(Cyan("src    - put your source files here."))
    fmt.Println(Cyan("lib    - libraries for the project go here."))
    if app.config.Tests {
        fmt.Println(Cyan("test   - put your files for unit testing here."))
    }
}

// Creates a template project that is ready to build and upload for application type
func (app App) createTemplateProject() (error) {
    return nil
}

// Prints all the commands relevant to application type
func (app App) printNextCommands() {
    fmt.Println(Cyan("`wio build -h`"))
    fmt.Println(Cyan("`wio run -h`"))
    fmt.Println(Cyan("`wio upload -h`"))

    if app.config.Tests {
        fmt.Println(Cyan("`wio build -h`"))
    }
}
