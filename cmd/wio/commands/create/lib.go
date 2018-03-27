// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates a library to be published
package create

import (
    "path/filepath"
    "os"
    "fmt"

    . "wio/cmd/wio/types"
    . "github.com/logrusorgru/aurora"
)

type Lib struct {
    config ConfigCreate
}

// Creates project structure for library type
func (lib Lib) createStructure() (error) {
    srcPath := lib.config.Directory + string(filepath.Separator) + "src"
    libPath := lib.config.Directory + string(filepath.Separator) + "lib"
    testPath := lib.config.Directory + string(filepath.Separator) + "test"
    wioPath := lib.config.Directory + string(filepath.Separator) + ".wio"

    err := os.MkdirAll(srcPath, os.ModePerm)
    if err != nil {
        return err
    }

    err = os.MkdirAll(libPath, os.ModePerm)
    if err != nil {
        return err
    }

    err = os.MkdirAll(wioPath, os.ModePerm)
    if err != nil {
        return err
    }

    err = os.MkdirAll(testPath, os.ModePerm)
    if err != nil {
        return err
    }

    return nil
}

// Prints the project structure for library type
func (lib Lib) printProjectStructure() {
    fmt.Println()
    fmt.Println(Cyan("src    - put your source files here."))
    fmt.Println(Cyan("lib    - libraries for the project go here."))
    fmt.Println(Cyan("test   - put your files for unit testing here."))
}

// Creates a template project that is ready to build and upload for library type
func (lib Lib) createTemplateProject() (error) {
    return nil
}

// Prints all the commands relevant to library type
func (lib Lib) printNextCommands() {
    fmt.Println(Cyan("`wio build -h`"))
    fmt.Println(Cyan("`wio run -h`"))
    fmt.Println(Cyan("`wio upload -h`"))
    fmt.Println(Cyan("`wio test -h`"))
}
