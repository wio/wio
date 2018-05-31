// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "github.com/urfave/cli"
    "path/filepath"
    "wio/cmd/wio/utils/io/log"
    "os"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "regexp"
    "wio/cmd/wio/utils"
    "bufio"
    "strings"
    "wio/cmd/wio/commands"
    "errors"
)

const (
    APP = "app"
    PKG = "pkg"
)

/// This structure wraps all the important features needed for a create and update command
type PacketCreate struct {
    directory string
    name      string
    board     string
    framework string
    platform  string
    ide       string
    tests     bool
}

type Create struct {
    Context *cli.Context
    Type    string
    Update  bool
    error   error
}

// get context for the command
func (create Create) GetContext() (*cli.Context) {
    return create.Context
}

// Executes the create command
func (create Create) Execute() {
    log.Norm.Cyan(false, "verifying project configuration ... ")
    commands.RecordError(create.error, "")

    if !create.Update && len(create.Context.Args()) < 2 {
        // When we are creating a project we need both directory and a board from args
        commands.RecordError(errors.New("project directory or Board not specified"), "failure")
    } else if !create.Update && len(create.Context.Args()) < 1 {
        // When we are updating a project we only need directory from args
        commands.RecordError(errors.New("project directory not specified"), "failure")
    }

    createPacket := &PacketCreate{}

    // fetch directory based on the argument
    directory, err := filepath.Abs(create.Context.Args()[0])
    commands.RecordError(err, "failure")

    // based on the command line context, get all the important fields
    createPacket.directory = directory
    createPacket.name = filepath.Base(directory)
    createPacket.framework = create.Context.String("framework")
    createPacket.platform = create.Context.String("platform")
    createPacket.ide = create.Context.String("ide")
    createPacket.tests = create.Context.Bool("tests")

    log.Norm.Green(true, "success")

    if create.Update {
        // if update command is called
        if create.checkUpdate(directory) {
            // after check if an update can be performed
            createPacket.board = create.Context.String("board")
            create.updateProject(createPacket)
        } else {
            // project structure is invalid or too distorted for an update to be applied
            message := `Based on your current project structure and wio.yml file, update cannot be performed.
You can use: wio create <app type> DIRECTORY BOARD`
            log.Norm.Cyan(true, message)
        }
    } else if create.createCheck(directory) {
        // if create command is called and it's check is passed
        createPacket.board = create.Context.Args()[1]

        create.createStructure(createPacket.directory, true)
        create.initialProjectSetup(createPacket)
        create.postPrint(createPacket, false)
    }
}

/// This method is a helper function that prints instructions to the console before "create"
/// or "update" command is executed. This has information about the project and it's ty[e
func (create Create) prePrint(createPacket *PacketCreate) {
    log.Norm.Yellow(false, "Project Name: ")
    log.Norm.Cyan(true, createPacket.name)
    log.Norm.Yellow(false, "Project Type: ")
    log.Norm.Cyan(true, create.Type)
    log.Norm.Yellow(false, "Project Path: ")
    log.Norm.Cyan(true, createPacket.directory)
}

/// This method is used to create the project. This erases everything that there is in the directory
/// and creates the project from scratch. This is the reason that this method is only called after
/// a check has been performed
func (create Create) createStructure(directory string, delete bool) {
    if delete {
        log.Norm.Cyan(false, "creating project structure ... ")
    } else {
        log.Norm.Cyan(false, "updating project structure ... ")
    }

    if delete && utils.PathExists(directory) {
        commands.RecordError(os.RemoveAll(directory), "failure")
    }

    commands.RecordError(os.MkdirAll(directory+io.Sep+"src", os.ModePerm), "failure")
    commands.RecordError(os.MkdirAll(directory+io.Sep+".wio"+io.Sep+"build", os.ModePerm), "failure")

    if create.Type == PKG {
        /// each package will have an include method
        commands.RecordError(os.MkdirAll(directory+io.Sep+"include", os.ModePerm),
            "failure")
    }

    if create.Type == PKG || create.Context.Bool("tests") {
        commands.RecordError(os.MkdirAll(directory+io.Sep+"tests", os.ModePerm), "failure")
    }

    log.Norm.Green(true, "success")
}

// This is a method that updates the created project. It will fix the structure and apply
// updates. It also makes sure all the configurations are applied
func (create Create) updateProject(createPacket *PacketCreate) {
    create.createStructure(createPacket.directory, false)
    create.updateProjectSetup(createPacket)
    create.postPrint(createPacket, true)
}

/// This is one of the most important step as this sets up the project when update command is used.
/// This also updates the wio.yml file so that it can be fixed and current configurations can be applied.
func (create Create) updateProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "updating project files ... ")

    commands.RecordError(copyTemplates(createPacket.directory, create.Type, createPacket.ide,
        "config"+io.Sep+"update_paths.json"), "failure")

    // get default configuration values
    defaults := types.DConfig{}
    commands.RecordError(io.AssetIO.ParseYml("config/defaults.yml", &defaults), "failure")

    // we have to copy README file
    var readmeStr string

    var config interface{}

    if create.Type == APP {
        projectConfig := &types.AppConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.directory+io.Sep+"wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.name
        // update the targets to make sure they are valid and there is a default target
        create.handleTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.Default_target].Board

        // check framework and platform
        checkFrameworkAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.name, projectConfig.MainTag.Platform, projectConfig.MainTag.Framework)
    } else {
        projectConfig := &types.PkgConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.directory+io.Sep+"wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.name
        // update the targets to make sure they are valid and there is a default target
        create.handleTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.Default_target].Board

        // check project version
        if projectConfig.MainTag.Version == "" {
            projectConfig.MainTag.Version = "0.0.1"
        }

        // make sure boards are updated in yml file
        if !utils.StringInSlice("ALL", projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = []string{"ALL"}
        } else if !utils.StringInSlice(createPacket.board, projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = append(projectConfig.MainTag.Board, createPacket.board)
        }

        // check frameworks and platform
        checkFrameworkArrayAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.name, projectConfig.MainTag.Platform,
            projectConfig.MainTag.Framework, projectConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfig(config, createPacket.directory+io.Sep+"wio.yml"),
        "failure")

    if utils.PathExists(createPacket.directory + io.Sep + "README.md") {
        content, err := io.NormalIO.ReadFile(createPacket.directory + io.Sep + "README.md")
        commands.RecordError(err, "failure")

        // only write if README is empty
        if string(content) == "" || string(content) == "\n" {
            // write readme file
            io.NormalIO.WriteFile(createPacket.directory+io.Sep+"README.md", []byte(readmeStr))
        }
    } else {
        // readme does not exist so write it
        io.NormalIO.WriteFile(createPacket.directory+io.Sep+"README.md", []byte(readmeStr))
    }

    log.Norm.Green(true, "success")
}

// This function checks if framework and platform are not empty. It future we can in force in valud
// frameworks and platforms using this
func checkFrameworkAndPlatform(framework *string, platform *string, defaults *types.DConfig) {
    if *framework == "" {
        *framework = defaults.Framework
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

// This function is similar to the above but in this case it checks if multiple frameworks are invalid
// and same goes for platform
func checkFrameworkArrayAndPlatform(framework *[]string, platform *string, defaults *types.DConfig) {
    if len(*framework) == 0 {
        *framework = append(*framework, defaults.Framework)
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

/// This is one of the most important step as this sets up the project when create command is used.
/// This also fills up the wio.yml file so that default configuration along with user choices
/// are applied.
func (create Create) initialProjectSetup(createPacket *PacketCreate) {
    log.Norm.Cyan(false, "creating template project ... ")

    defaultTarget := "main"

    commands.RecordError(copyTemplates(createPacket.directory, create.Type, createPacket.ide,
        "config"+io.Sep+"create_paths.json"), "failure")

    // we have to copy README file
    var readmeStr string

    var config interface{}

    if create.Type == APP {
        projectConfig := &types.AppConfig{}

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.ide
        projectConfig.MainTag.Platform = createPacket.platform
        projectConfig.MainTag.Framework = createPacket.framework
        projectConfig.MainTag.Name = createPacket.name
        projectConfig.TargetsTag.Default_target = defaultTarget
        targets := make(map[string]*types.TargetTag, 1)
        projectConfig.TargetsTag.Targets = targets

        targetsTag := projectConfig.TargetsTag
        create.handleTargets(&targetsTag, createPacket.board)

        config = projectConfig

        // update readme app content
        readmeStr = getReadmeApp(createPacket.name, projectConfig.MainTag.Platform, projectConfig.MainTag.Framework)
    } else {
        projectConfig := &types.PkgConfig{}
        defaultTarget = "tests"

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.ide
        projectConfig.MainTag.Platform = createPacket.platform
        projectConfig.MainTag.Framework = []string{createPacket.framework}
        projectConfig.MainTag.Name = createPacket.name
        projectConfig.MainTag.Version = "0.0.1"
        projectConfig.TargetsTag.Default_target = defaultTarget
        targets := make(map[string]*types.TargetTag, 0)
        projectConfig.TargetsTag.Targets = targets
        projectConfig.MainTag.Board = []string{createPacket.board}

        targetsTag := projectConfig.TargetsTag
        create.handleTargets(&targetsTag, createPacket.board)

        config = projectConfig

        // update readme pkg content
        readmeStr = getReadmePkg(createPacket.name, projectConfig.MainTag.Platform,
            projectConfig.MainTag.Framework, projectConfig.MainTag.Board)
    }

    commands.RecordError(utils.PrettyPrintConfig(config, createPacket.directory+io.Sep+"wio.yml"),
        "failure")

    // write readme file
    io.NormalIO.WriteFile(createPacket.directory+io.Sep+"README.md", []byte(readmeStr))

    log.Norm.Green(true, "success")
}

// returns README string for wio application
func getReadmeApp(name string, platform string, framework string) (string) {
    readmeContent, err := io.AssetIO.ReadFile("templates/readme/APP_README.md")
    commands.RecordError(err, "failure")

    readmeStr := strings.Replace(string(readmeContent), "{{PROJECT_NAME}}", name, -1)
    readmeStr = strings.Replace(readmeStr, "{{PLATFORM}}", strings.ToUpper(platform), -1)
    readmeStr = strings.Replace(readmeStr, "{{FRAMEWORK}}", framework, -1)

    return readmeStr
}

// returns README string for wio package
func getReadmePkg(name string, platform string, framework []string, board []string) (string) {
    readmeContent, err := io.AssetIO.ReadFile("templates/readme/PKG_README.md")
    commands.RecordError(err, "failure")

    readmeStr := strings.Replace(string(readmeContent), "{{PROJECT_NAME}}", name, -1)
    readmeStr = strings.Replace(readmeStr, "{{PLATFORM}}", strings.ToUpper(platform), -1)
    readmeStr = strings.Replace(readmeStr, "{{FRAMEWORKS}}",
        strings.Join(framework, ","), -1)
    readmeStr = strings.Replace(readmeStr, "{{BOARDS}}",
        strings.Join(board, ","), -1)

    return readmeStr
}

/// This method handles the targets that a user can create and what these targets are
/// in wio.yml file. It targets are not there, it will create a default target. Unless
/// it will keep the targets that are already there
func (create Create) handleTargets(targetsTag *types.TargetsTag, board string) {
    defaultTarget := &types.TargetTag{}

    if targetsTag.Targets == nil {
        targetsTag.Targets = make(map[string]*types.TargetTag)
    }

    if target, ok := targetsTag.Targets[targetsTag.Default_target]; ok {
        defaultTarget.Board = target.Board
        defaultTarget.Compile_flags = target.Compile_flags
        targetsTag.Targets[targetsTag.Default_target] = defaultTarget
    } else {
        defaultTarget.Board = board
        targetsTag.Targets[targetsTag.Default_target] = defaultTarget
    }
}

/// This method prints next steps for any type of create/update command. This will help user
/// decide what they can do next
func (create Create) postPrint(createPacket *PacketCreate, update bool) {
    if !update {
        log.Norm.Yellow(true, "Project Creation finished!! Summary: ")
    } else {
        log.Norm.Yellow(true, "Project Update finished!! Summary: ")
    }

    log.Norm.Cyan(false, "project name: ")
    log.Norm.Cyan(true, createPacket.name)
    log.Norm.Cyan(false, "project type: ")
    log.Norm.Cyan(true, create.Type)
    log.Norm.Cyan(false, "project path: ")
    log.Norm.Cyan(true, createPacket.directory)

    log.Norm.Yellow(true, "Check following commands:")

    log.Norm.Cyan(true, "`wio build -h`")
    log.Norm.Cyan(true, "`wio run -h`")
    log.Norm.Cyan(true, "`wio pac -h`")

    if create.Context.Bool("tests") {
        log.Norm.Cyan(true, "`wio test -h`")
    }
}

/// This method is a crucial peace of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func (create Create) createCheck(directory string) (bool) {
    if !utils.PathExists(directory) {
        return true
    }

    if status, err := utils.IsEmpty(directory); err != nil {
        commands.RecordError(err, "")
        return false
    } else if status {
        return true
    } else {
        message := `The directory is not empty!!
This action will erase everything and will create a new project.
An alternative is to do: wio update <app type> DIRECTORY
Please type y/yes to indicate creation and anything else to indicate abortion: `
        log.Norm.Cyan(false, message)
        reader := bufio.NewReader(os.Stdin)
        text, err := reader.ReadString('\n')
        commands.RecordError(err, "")

        text = strings.TrimSuffix(strings.ToLower(text), "\n")

        if text == "y" || text == "yes" {
            log.Norm.Write(true, "")
            return true
        } else {
            return false
        }
    }
}

/// This method checks if an update can be performed. This basically checks for compatibility issues
/// and configurations to make sure the update can be performed
func (create Create) checkUpdate(directory string) (bool) {
    wioPath := directory + io.Sep + "wio.yml"

    if status, err := utils.IsEmpty(directory); err != nil {
        // if there is an error
        commands.RecordError(err, "")
        return false
    } else if status {
        // if the folder is empty
        return false
    } else if !utils.PathExists(wioPath) {
        // if wio.yml file does not exist
        return false
    } else {
        if create.Type == APP {
            config := &types.AppConfig{}
            if err := io.NormalIO.ParseYml(wioPath, config); err != nil {
                // can't parse wio.yml file for app type

                // check if it contains "app:" tag
                if appContent, err := io.NormalIO.ReadFile(wioPath); err != nil {
                    return false
                } else {
                    // if it does make a new one and return true
                    if strings.Contains(string(appContent), "app:") {
                        commands.RecordError(os.Remove(wioPath), "")
                        return true
                    }
                }

                return false
            } else if config.MainTag.Name == "" {
                // type of the project is wrong compared to the one specified
                return false
            } else {
                // all checks passed
                return true
            }
        } else {
            config := &types.PkgConfig{}
            if err := io.NormalIO.ParseYml(wioPath, config); err != nil {
                // can't parse wio.yml file for pkg type

                // check if it contains "pkg:" tag
                if pkgContent, err := io.NormalIO.ReadFile(wioPath); err != nil {
                    return false
                } else {
                    // if it does make a new one and return true
                    if strings.Contains(string(pkgContent), "pkg:") {
                        commands.RecordError(os.Remove(wioPath), "")
                        return true
                    }
                }

                return false
            } else if config.MainTag.Name == "" {
                // type of the project is wrong compared to the one specified
                return false
            } else {
                // all checks passed
                return true
            }
        }
    }
}

/// This function copies the templates needed to set up the project. It uses parsePathsAndCopy
/// function to parse the paths.json file and then get files based on the project type.
func copyTemplates(projectPath string, appType string, ide string, jsonPath string) (error) {
    strArray := make([]string, 0)
    strArray = append(strArray, appType+"-gen")

    if ide == "clion" {
        strArray = append(strArray, appType+"-clion")
    }

    if err := parsePathsAndCopy(jsonPath, projectPath, strArray); err != nil {
        return err
    }

    return nil
}

/// This function Parses the paths.json file and uses that to get the required files to set
/// up the wio project. This also decides the override based on the paths.json file
func parsePathsAndCopy(jsonPath string, projectPath string, tags []string) (error) {
    var paths = Paths{}
    if err := io.AssetIO.ParseJson(jsonPath, &paths); err != nil {
        return err
    }

    var re, e = regexp.Compile(`{{.+}}`)
    if e != nil {
        return e
    }

    var sources []string
    var destinations []string
    var overrides []bool

    for i := 0; i < len(paths.Paths); i++ {
        for t := 0; t < len(tags); t++ {
            if paths.Paths[i].Id == tags[t] {
                sources = append(sources, paths.Paths[i].Src)

                destination, e := filepath.Abs(re.ReplaceAllString(paths.Paths[i].Des, projectPath))
                if e != nil {
                    return e
                }

                destinations = append(destinations, destination)
                overrides = append(overrides, paths.Paths[i].Override)
            }
        }
    }

    return io.AssetIO.CopyMultipleFiles(sources, destinations, overrides)
}
