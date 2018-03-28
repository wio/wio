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
    "github.com/spf13/viper"
)

type App struct {
    config ConfigCreate
    ioData IOData
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
    strArray := make([]string, 1)
    strArray[0] = "app-gen"

    if app.config.Ide == "clion" {
        strArray = append(strArray, "app-clion")
    }

    path := "assets" + sep + "config" + sep + "paths.json"

    paths, err := ParsePathsAndCopy(path, app.config, app.ioData, strArray)
    if err != nil {
        return err
    }

    // fill the templates
    return app.fillTemplates(paths)

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

func (app App) fillTemplates(paths map[string]string) (error) {
    err := app.FillConfig(paths)
    if err != nil {
        return err
    }

    return nil
}

// Handles config file for app
func (app App) FillConfig(paths map[string]string) (error) {
    viper.SetConfigName("project")
    viper.AddConfigPath(app.config.Directory)
    err := viper.ReadInConfig()

    if err != nil {
        return err
    }

    // parse app
    wioApp, err := getAppStruct(viper.GetStringMap("app"))
    if err != nil {
        return err
    }

    // parse libraries
    libraries := getLibrariesStruct(viper.GetStringMap("libraries"))

    // write new config
    wioApp.Ide = app.config.Ide
    wioApp.Platform = app.config.Platform
    wioApp.Framework = app.config.Framework
    wioApp.Name = filepath.Base(app.config.Directory)

    if wioApp.Default_target == "" {
        wioApp.Default_target = "main"
    }

    wioApp.Targets[wioApp.Default_target] = &TargetStruct{}
    wioApp.Targets[wioApp.Default_target].Board = app.config.Board

    viper.Set("app", wioApp)
    viper.Set("libraries", libraries)

    configData := viper.AllSettings()
    infoPath := "assets" + sep + "templates" + sep + "config" + sep + "project-app-help"
    if err != nil {
        return err
    }

    err = PrettyWriteConfig(infoPath, configData, app.ioData, paths["project.yml"])
    if err != nil {
        return err
    }

    return nil
}

func (app App) FillCMake(paths map[string]string) (error) {
    return nil
}
