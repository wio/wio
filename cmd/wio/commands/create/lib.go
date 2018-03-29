// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates a library to be published
package create

import (
    "path/filepath"
    "os"

    . "wio/cmd/wio/io"
    . "wio/cmd/wio/types"
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

    err = os.MkdirAll(testPath, os.ModePerm)
    if err != nil {
        return err
    }
    Verb.Verbose(`* Created "test" folder` + "\n")

    return nil
}

// Prints the project structure for library type
func (lib Lib) printProjectStructure() {
    Norm.Cyan("src    - put your source files here.\n")
    Norm.Cyan("lib    - libraries for the project go here.\n")
    Norm.Cyan("test   - put your files for unit testing here.\n")
}

// Creates a template project that is ready to build and upload for library type
func (lib Lib) createTemplateProject() (error) {
    strArray := make([]string, 1)
    strArray[0] = "lib-gen"

    Verb.Verbose("\n")
    if lib.config.Ide == "clion" {
        Verb.Verbose("* Clion Ide available so ide template set up will be used\n")
        strArray = append(strArray, "lib-clion")
    } else {
        Verb.Verbose("* General template setup will be used\n")
    }

    path := "assets" + sep + "config" + sep + "paths.json"
    paths, err := ParsePathsAndCopy(path, lib.config, lib.ioData, strArray)
    if err != nil {
        return err
    }
    Verb.Verbose("* All Template files created in their right position\n")

    // fill the templates
    return lib.fillTemplates(paths)
}

// Prints all the commands relevant to library type
func (lib Lib) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")
    Norm.Cyan("`wio test -h`\n")
}

func (lib Lib) fillTemplates(paths map[string]string) (error) {
    err := lib.FillConfig(paths)
    if err != nil {
        return err
    }
    Verb.Verbose("* Finished Filling/Updating Project.yml template\n")

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
    Verb.Verbose("* Loaded Project.yml file template\n")

    // parse app
    wioLib, err := getLibStruct(viper.GetStringMap("lib"))
    if err != nil {
        return err
    }
    Verb.Verbose("* Parsed lib tag of the template\n")

    // parse libraries
    libraries := getLibrariesStruct(viper.GetStringMap("libraries"))
    Verb.Verbose("* Parsed libraries tag of the template\n")

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

    err = PrettyWriteConfig(infoPath, configData, lib.ioData, paths["project.yml"])
    if err != nil {
        return err
    }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return nil
}

func (lib Lib) FillCMake(paths map[string]string) (error) {
    return nil
}
