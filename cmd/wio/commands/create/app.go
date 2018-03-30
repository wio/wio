// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates an executable application
package create

import (
    "path/filepath"

    . "wio/cmd/wio/utils/io"
    . "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/types"
)


// Creates project structure for application type
func (app App) createStructure() (error) {
    Verb.Verbose("\n")
    if err := createStructure(app.args.Directory, "src", "lib", ".wio/targets"); err != nil {
        return err
    }

    if app.args.Tests {
        if err := createStructure(app.args.Directory, "test"); err != nil {
            return err
        }
    }

    return nil
}

// Prints the project structure for application type
func (app App) printProjectStructure() {
    Norm.Cyan("src    - put your source files here.\n")
    Norm.Cyan("lib    - libraries for the project go here.\n")
    if app.args.Tests {
        Norm.Cyan("test   - put your files for unit testing here.\n")
    }
}

// Creates a template project that is ready to build and upload for application type
func (app App) createTemplateProject() (error) {
    config := &types.AppConfig{}
    var err error

    if err = copyTemplates(app.args); err != nil { return err }
    if config, err = app.FillConfig(); err != nil { return err }

    // create cmake files for each target
    copyTargetCMakes(app.args.Directory, app.args.AppType, config.MainTag.Targets)

    return nil
}

// Prints all the commands relevant to application type
func (app App) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")

    if app.args.Tests {
        Norm.Cyan("`wio test -h`\n")
    }
}

// Handles config file for app
func (app App) FillConfig() (*types.AppConfig, error) {
    Verb.Verbose("* Loaded Project.yml file template\n")

    appConfig := types.AppConfig{}
    if err := NormalIO.ParseYml(app.args.Directory + Sep + "project.yml", &appConfig);
    err != nil { return nil, err }

    // make modifications to the data
    appConfig.MainTag.Ide = app.args.Ide
    appConfig.MainTag.Platform = app.args.Platform
    appConfig.MainTag.Framework = app.args.Framework
    appConfig.MainTag.Name = filepath.Base(app.args.Directory)

    if appConfig.MainTag.Default_target == "" {
        appConfig.MainTag.Default_target = "main"
    }

    appConfig.MainTag.Targets[appConfig.MainTag.Default_target] = &types.TargetSubTags{}
    appConfig.MainTag.Targets[appConfig.MainTag.Default_target].Board = app.args.Board

    Verb.Verbose("* Modified information in the configuration\n")

    if err := PrettyPrintConfig(app.args.AppType, &appConfig, app.args.Directory + Sep + "project.yml");
    err != nil { return nil, err }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return &appConfig, nil
}

func (app App) FillCMake(paths map[string]string) (error) {
    return nil
}
