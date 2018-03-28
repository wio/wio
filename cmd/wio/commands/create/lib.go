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
    "github.com/spf13/viper"
    "wio/cmd/wio/utils"
)

type Lib struct {
    config ConfigCreate
    ioData IOData
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
    strArray := make([]string, 1)
    strArray[0] = "lib-gen"

    if lib.config.Ide == "clion" {
        strArray = append(strArray, "lib-clion")
    }

    path := "assets" + sep + "config" + sep + "paths.json"
    paths, err := ParsePathsAndCopy(path, lib.config, lib.ioData, strArray)
    if err != nil {
        return err
    }

    // fill the templates
    return lib.fillTemplates(paths)

    return nil
}

// Prints all the commands relevant to library type
func (lib Lib) printNextCommands() {
    fmt.Println(Cyan("`wio build -h`"))
    fmt.Println(Cyan("`wio run -h`"))
    fmt.Println(Cyan("`wio upload -h`"))
    fmt.Println(Cyan("`wio test -h`"))
}

func (lib Lib) fillTemplates(paths map[string]string) (error) {
    err := lib.FillConfig(paths)
    if err != nil {
        return err
    }

    return nil
}

// Handles config file for lib
func (lib Lib) FillConfig(paths map[string]string) (error) {
    viper.SetConfigName("project")
    viper.AddConfigPath(lib.config.Directory)
    err := viper.ReadInConfig()

    if err != nil {
        return err
    }

    // parse app
    wioLib, err := getLibStruct(viper.GetStringMap("lib"))
    if err != nil {
        return err
    }

    // parse libraries
    libraries := getLibrariesStruct(viper.GetStringMap("libraries"))

    // write new config
    wioLib.Ide = lib.config.Ide
    wioLib.Platform = lib.config.Platform
    wioLib.Framework = utils.AppendIfMissing(wioLib.Framework, lib.config.Framework)
    wioLib.Name = filepath.Base(lib.config.Directory)

    if wioLib.Default_target == "" {
        wioLib.Default_target = "test"
    }

    wioLib.Targets[wioLib.Default_target] = &TargetStruct{}
    wioLib.Targets[wioLib.Default_target].Board = lib.config.Board

    viper.Set("lib", wioLib)
    viper.Set("libraries", libraries)

    configData := viper.AllSettings()
    infoPath := "assets" + sep + "templates" + sep + "config" + sep + "project-lib-help"
    if err != nil {
        return err
    }

    err = PrettyWriteConfig(infoPath, configData, lib.ioData, paths["project.yml"])
    if err != nil {
        return err
    }

    return nil
}

func (lib Lib) FillCMake(paths map[string]string) (error) {
    return nil
}
