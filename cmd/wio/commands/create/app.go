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
    "wio/cmd/wio/utils"
)

type App struct {
    config ConfigCreate
}

// Creates project structure for application type
func (app App) createStructure() (error) {
    srcPath := app.config.Directory + string(filepath.Separator) + "src"
    libPath := app.config.Directory + string(filepath.Separator) + "lib"
    wioPath := app.config.Directory + string(filepath.Separator) + ".wio"

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
    root, err := utils.GetExecutableRootPath()
    rootPath = root

    if err != nil {
        return err
    }

    // copy config file
    configPathSrc := GetTemplatesRelativeFile("config" + sep + "project-app.yml")
    configPathDest := app.config.Directory + sep + "project.yml"

    err = utils.Copy(configPathSrc, configPathDest, false)

    if err != nil {
        return err
    }

    // copy main.cpp from template in src folder
    mainPathSrc := GetTemplatesRelativeFile("sample-program" + sep + "main.cpp")
    mainPathDest := app.config.Directory + sep + "src" + sep + "main.cpp"

    err = utils.Copy(mainPathSrc, mainPathDest, false)

    if err != nil {
        return err
    }

    // copy CMakeLists.txt from template in .wio folder
    cmakeSrc := GetTemplatesRelativeFile( "cmake" + sep + "CMakeLists.txt.tpl")
    cmakeDest := app.config.Directory + sep + ".wio" + sep + "CMakeLists.txt"

    err = utils.Copy(cmakeSrc, cmakeDest, true)

    if err != nil {
        return err
    }

    // copy .gitignore file to the root
    gitignoreSrc := GetTemplatesRelativeFile("gitignore" + sep + ".gitignore-general")

    if app.config.Ide == "clion" {
        gitignoreSrc = GetTemplatesRelativeFile("gitignore" + sep + ".gitignore-clion")
    }
    gitignoreDest := app.config.Directory + sep + ".gitignore"

    err = utils.Copy(gitignoreSrc, gitignoreDest, false)

    if err != nil {
        return err
    }

    // fill the templates
    return fillTemplates()

    return nil
}

// Prints all the commands relevant to application type
func (app App) printNextCommands() {
    fmt.Println(Cyan("`wio build -h`"))
    fmt.Println(Cyan("`wio run -h`"))
    fmt.Println(Cyan("`wio upload -h`"))

    if app.config.Tests {
        fmt.Println(Cyan("`wio test -h`"))
    }
}

func fillTemplates() (error) {
    return nil
}
