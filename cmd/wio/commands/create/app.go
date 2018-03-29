// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates an executable application
package create

import (
    "os"
    "path/filepath"

    . "wio/cmd/wio/io"
    . "wio/cmd/wio/types"
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
    Verb.Verbose("\n")
    Verb.Verbose(`* Created "src" folder` + "\n")

    err = os.MkdirAll(libPath, os.ModePerm)
    if err != nil {
        return err
    }
    Verb.Verbose(`* Created "lib" folder` + "\n")

    err = os.MkdirAll(wioPath, os.ModePerm)
    if err != nil {
        return err
    }
    Verb.Verbose(`* Created ".wio" folder` + "\n")

    if app.config.Tests {
        testPath := app.config.Directory + string(filepath.Separator) + "test"
        err = os.MkdirAll(testPath, os.ModePerm)
        if err != nil {
            return err
        }
        Verb.Verbose(`* Created "tests" folder` + "\n")
    }

    return nil
}

// Prints the project structure for application type
func (app App) printProjectStructure() {
    Norm.Cyan("src    - put your source files here.\n")
    Norm.Cyan("lib    - libraries for the project go here.\n")
    if app.config.Tests {
        Norm.Cyan("test   - put your files for unit testing here.\n")
    }
}

// Creates a template project that is ready to build and upload for application type
func (app App) createTemplateProject() (error) {
    strArray := make([]string, 1)
    strArray[0] = "app-gen"

    Verb.Verbose("\n")
    if app.config.Ide == "clion" {
        Verb.Verbose("* Clion Ide available so ide template set up will be used\n")
        strArray = append(strArray, "app-clion")
    } else {
        Verb.Verbose("* General template setup will be used\n")
    }

    path := "assets" + sep + "config" + sep + "paths.json"
    paths, err := ParsePathsAndCopy(path, app.config, app.ioData, strArray)
    if err != nil {
        return err
    }
    Verb.Verbose("* All Template files created in their right position\n")

    // fill the templates
    return app.fillTemplates(paths)
}

// Prints all the commands relevant to application type
func (app App) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")

    if app.config.Tests {
        Norm.Cyan("`wio test -h`\n")
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
    Verb.Verbose("* Loaded Project.yml file template\n")

    // parse app
    wioApp, err := getAppStruct(viper.GetStringMap("app"))
    if err != nil {
        return err
    }
    Verb.Verbose("* Parsed app tag of the template\n")

    // parse libraries
    libraries := getLibrariesStruct(viper.GetStringMap("libraries"))
    Verb.Verbose("* Parsed libraries tag of the template\n")

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

    err = PrettyWriteConfig(infoPath, configData, app.ioData, paths["project.yml"])
    if err != nil {
        return err
    }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return nil
}

func (app App) FillCMake(paths map[string]string) (error) {
    return nil
}
