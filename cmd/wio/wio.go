// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
 package main

import (
<<<<<<< HEAD
    "time"
    "os"
    "log"
    "fmt"
    "path/filepath"

    "github.com/urfave/cli"
    commandCreate "wio/cmd/wio/commands/create"
    . "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils/io"
)

//go:generate go-bindata -nomemcopy -prefix ../../ ../../assets/config/... ../../assets/templates/...
func main()  {
    // override help template
    cli.AppHelpTemplate =
`Wio a simplified development process for embedded applications.
=======
    "wio/cmd/wio/utils/io"
    "github.com/urfave/cli"
    "time"
    "os"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/pac"
    "wio/cmd/wio/utils/io/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/commands/build"
    "wio/cmd/wio/commands/clean"
    "wio/cmd/wio/commands/run"
)

func main() {
    // override help template
    cli.AppHelpTemplate =
        `Wio a simplified development process for embedded applications.
>>>>>>> More commands and minor fixes (#37)
Create, Build, Test, and Upload AVR projects from Commandline.

Common Commands:
    
    wio create <app type> [options] <output directory>
        Create a new Wio project in the specified directory.
    
    wio build [options]
        Build the Wio project based on all the configurations defined
    
    wio upload [options]
        Upload the Wio project to an attached embedded device
    
    wio run [options]
        Builds, Tests, and Uploads the Wio projects

Usage: {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
Global options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
Available commands:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
Global options:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

Copyright:
   {{.Copyright}}
   {{end}}{{if .Version}}
Vesrion:
   {{.Version}}
   {{end}}
Run "wio command <help>" for more information about a command.
`

<<<<<<< HEAD
cli.CommandHelpTemplate =
`{{.Usage}}
=======
    cli.CommandHelpTemplate =
        `{{.Usage}}
>>>>>>> More commands and minor fixes (#37)

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

Category:
   {{.Category}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .VisibleFlags}}

Available commands:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`

<<<<<<< HEAD
cli.SubcommandHelpTemplate =
`{{.Usage}}
=======
    cli.SubcommandHelpTemplate =
        `{{.Usage}}
>>>>>>> More commands and minor fixes (#37)

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

Available commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`
    // get default configuration values
<<<<<<< HEAD
    defaults := DConfig{}
    err := io.AssetIO.ParseYml("config/defaults.yml", &defaults)

    if err != nil {
        log.Fatal(err)
    }

=======
    defaults := types.DConfig{}
    err := io.AssetIO.ParseYml("config/defaults.yml", &defaults)
    if err != nil {
        log.Error(false, err.Error())
    }

    // command that will be executed
    var command commands.Command

>>>>>>> More commands and minor fixes (#37)
    app := cli.NewApp()
    app.Name = "wio"
    app.Version = defaults.Version
    app.EnableBashCompletion = true
    app.Compiled = time.Now()
    app.Copyright = "Copyright (c) 2018 Waterloop"
    app.Usage = "Create, Build and Upload AVR projects"

<<<<<<< HEAD
    app.Flags = []cli.Flag {
        cli.BoolFlag{Name: "verbose",
            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
            },
    }

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project. Also acts as updater/fixer",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "lib",
                    Usage:     "Creates a wio library, intended to be used by other people",
                    UsageText: "wio create lib <DIRECTORY> <BOARD> [command options]",
=======
    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create pkg <DIRECTORY> <BOARD> [command options]",
>>>>>>> More commands and minor fixes (#37)
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
<<<<<<< HEAD
                    },
                    Action: func(c *cli.Context) error {
                        // check if user defined a board
                        if len(c.Args()) == 0 {
                            fmt.Println("A Directory path/name is needed to create this wio library!")
                            fmt.Println("\nExecute `wio create app -h` for more details and help")
                            os.Exit(1)
                        } else if len(c.Args()) == 1 {
                            fmt.Println("A Board is needed to create this wio library!")
                            fmt.Println("\nExecute `wio create app -h` for more details and help")
                            os.Exit(1)
                        }

                        directory, _ := filepath.Abs(c.Args()[0])

                        libArgs := CliArgs{
                            AppType: "lib",
                            Directory: directory,
                            Board: c.Args()[1],
                            Framework: c.String("framework"),
                            Platform: c.String("platform"),
                            Ide: c.String("ide"),
                            Tests: true,
                        }
                        turnVerbose(c.GlobalBool("verbose"))

                        commandCreate.Execute(libArgs)

                        return nil
=======
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: false}
>>>>>>> More commands and minor fixes (#37)
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
<<<<<<< HEAD
                    },
                    Action: func(c *cli.Context) error {
                        // check if user defined a board
                        if len(c.Args()) == 0 {
                            fmt.Println("A Directory path/name is needed to create this wio application!")
                            fmt.Println("\nExecute `wio create app -h` for more details and help")
                            os.Exit(1)
                        } else if len(c.Args()) == 1 {
                            fmt.Println("A Board is needed to create this wio application!")
                            fmt.Println("\nExecute `wio create app -h` for more details and help")
                            os.Exit(1)
                        }

                        directory, _ := filepath.Abs(c.Args()[0])

                        appArgs := CliArgs{
                            AppType: "app",
                            Directory: directory,
                            Board: c.Args()[1],
                            Framework: c.String("framework"),
                            Platform: c.String("platform"),
                            Ide: c.String("ide"),
                            Tests: c.Bool("tests"),
                        }
                        turnVerbose(c.GlobalBool("verbose"))

                        commandCreate.Execute(appArgs)

                        return nil
=======
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: false}
                    },
                },
            },
        },
        {
            Name:  "update",
            Usage: "Updates the current project and fixes any issues.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Updates a wio package, intended to be used by other people",
                    UsageText: "wio update pkg <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: true}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Updates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio update app <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: true}
>>>>>>> More commands and minor fixes (#37)
                    },
                },
            },
        },
        {
            Name:      "build",
<<<<<<< HEAD
            Usage:     "Builds the project",
=======
            Usage:     "Builds the wio project.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
<<<<<<< HEAD
                    Usage: "Build a specified target instead of building all the targets",
                    Value: defaults.Btarget,
                },
            },
            Action: func(c *cli.Context) error {
                /*
                build based on the type of project (from config file)
                 */
                return nil
=======
                    Usage: "Build a specified target instead of building the default",
                    Value: defaults.Btarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = build.Build{Context: c}
>>>>>>> More commands and minor fixes (#37)
            },
        },
        {
            Name:      "clean",
<<<<<<< HEAD
            Usage:     "Cleans all the build files for the project",
            UsageText: "wio clean",
            Action: func(c *cli.Context) error {
                return nil
=======
            Usage:     "Cleans all the build files for the project.",
            UsageText: "wio clean",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "target",
                    Usage: "Cleans build files for a specified target instead of cleaning all the targets",
                    Value: defaults.Btarget,
                },
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = clean.Clean{Context: c}
>>>>>>> More commands and minor fixes (#37)
            },
        },
        {
<<<<<<< HEAD
            Name:      "upload",
<<<<<<< HEAD
            Usage:     "Uploads the project to a device",
=======
            Usage:     "Uploads the project to a device.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio upload [command options]",
            Flags: []cli.Flag{
<<<<<<< HEAD
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
                },
<<<<<<< HEAD
=======
=======
>>>>>>> adding upload support
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
>>>>>>> More commands and minor fixes (#37)
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
<<<<<<< HEAD
=======
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
>>>>>>> More commands and minor fixes (#37)
            },
            Action: func(c *cli.Context)  {
                command = upload.Upload{Context: c}
            },
        },
        {
            Name:      "run",
<<<<<<< HEAD
            Usage:     "Builds, Tests, and Uploads the project to a device",
=======
            Usage:     "Builds and Uploads the project to a device.",
>>>>>>> More commands and minor fixes (#37)
=======
            Name:      "run",
            Usage:     "Builds and Uploads the project to a device. \n" +
                "In order to trigger upload specify port flag.",
>>>>>>> build, upload, clean and run commands finished
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
<<<<<<< HEAD
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
=======
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
>>>>>>> More commands and minor fixes (#37)
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
<<<<<<< HEAD
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
=======
                cli.StringFlag{Name: "dir",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
            },
            Action: func(c *cli.Context) {
                command = run.Run{Context: c}
>>>>>>> More commands and minor fixes (#37)
            },
        },
        {
            Name:      "test",
<<<<<<< HEAD
            Usage:     "Runs unit tests available in the project",
=======
            Usage:     "Runs unit tests available in the project.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio test",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
<<<<<<< HEAD
=======
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
>>>>>>> More commands and minor fixes (#37)
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "monitor",
<<<<<<< HEAD
            Usage:     "Runs the serial monitor",
=======
            Usage:     "Runs the serial monitor.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio monitor [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "gui",
                    Usage: "Runs the GUI version of the serial monitor tool",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
<<<<<<< HEAD
=======
                cli.BoolFlag{Name: "verbose",
                    Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                },
>>>>>>> More commands and minor fixes (#37)
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doctor",
<<<<<<< HEAD
            Usage:     "Show information about the installed tooling",
=======
            Usage:     "Guide development tooling and system configurations.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio doctor",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
<<<<<<< HEAD
            Name:      "configure",
            Usage:     "Configures paths for the tools used for development",
            UsageText: "wio configure [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "arduino-sdk-dir",
                    Usage: "path to Arduino SDK",
                },
                cli.StringFlag{Name: "make-path",
                    Usage: "path to `make` tool",
                },
                cli.StringFlag{Name: "cmake-path",
                    Usage: "Path to `cmake` tool",
                },
                cli.StringFlag{Name: "avr-path",
                    Usage: "Path to AVR libraries",
                },
                cli.StringFlag{Name: "arm-path",
                    Usage: "Path to ARM libraries",
                },
            },
            Action: func(c *cli.Context) error {
                /*
                If no flag provided, show current settings
                 */
                return nil
            },
        },
        {
            Name:      "analyze",
            Usage:     "Analyzes C/C++ code statically",
=======
            Name:      "analyze",
            Usage:     "Analyzes C/C++ code statically.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio analyze",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doxygen",
<<<<<<< HEAD
            Usage:     "Runs doxygen tool to create documentation for the code",
=======
            Usage:     "Runs doxygen tool to create documentation for the code.",
>>>>>>> More commands and minor fixes (#37)
            UsageText: "wio doxygen",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
<<<<<<< HEAD
            Name:  "packager",
            Usage: "Package manager for Wio projects",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:  "get",
                    Usage: "Gets all the packages being used in the project",
=======
            Name:  "pac",
            Usage: "Package manager for Wio projects.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "get",
                    Usage:     "Gets all the libraries mentioned in wio.yml file and vendor folder",
                    UsageText: "wio libraries get [command options]",
>>>>>>> More commands and minor fixes (#37)
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "clean",
                            Usage: "Cleans all the current packages and re get all of them",
                        },
<<<<<<< HEAD
                    },
                    Action: func(c *cli.Context) error {
                        return nil
                    },
                },
                cli.Command{
                    Name:  "update",
                    Usage: "Updates all the packages being used in the project and makes sure they are correct version",
                    Action: func(c *cli.Context) error {
                        return nil
                    },
                },
            },
        },
        {
            Name:  "tool",
            Usage: "Contains various tools related to setup, initialize and upgrade of Wio",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "setup",
                    Usage:     "When tool is newly installed, it sets up the tool for the machine",
                    UsageText: "wio setup",
                    Action: func(c *cli.Context) error {
                        return nil
                    },
                },
                cli.Command{
                    Name:      "upgrade",
                    Usage:     "Upgrades the current version of the program",
                    UsageText: "wio upgrade [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "version",
                            Usage: "Specify the exact version to upgrade/downgrade wio to",
                            Value: defaults.Version,
                        },
                    },
                    Action: func(c *cli.Context) error {
                        return nil
=======
                        cli.StringFlag{Name: "version_control",
                            Usage: "Specify the version control tool to usage",
                            Value: "git",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.GET}
                    },
                },
                cli.Command{
                    Name:      "update",
                    Usage:     "Updates all the libraries mentioned in wio.yml file and vendor folder.",
                    UsageText: "wio libraries update [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "version_control",
                            Usage: "Specify the version control tool to usage",
                            Value: "git",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.UPDATE}
                    },
                },
                cli.Command{
                    Name:      "collect",
                    Usage:     "Creates vendor folder and puts all the libraries in that folder",
                    UsageText: "wio libraries collect [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "path",
                            Usage: "Path to collect a library instead of collecting all of them",
                            Value: "none",
                        },
                        cli.BoolFlag{Name: "verbose",
                            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = pac.Pac{Context: c, Type: pac.COLLECT}
>>>>>>> More commands and minor fixes (#37)
                    },
                },
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

<<<<<<< HEAD
    err = app.Run(os.Args)

    if err != nil {
        log.Fatal(err)
    }
}

// Set's verbose mode on
func turnVerbose(value bool) {
    if value == true {
        io.SetVerbose()
    }
=======
    if err = app.Run(os.Args); err != nil {
        panic(err)
    }

    // execute the command
    if command != nil {
        // check if verbose flag is true
        if command.GetContext().Bool("verbose") {
            log.SetVerbose()
        }

        command.Execute()
    }
}

func getCurrDir() (string) {
    directory, err := os.Getwd()
    commands.RecordError(err, "")
    return directory
>>>>>>> More commands and minor fixes (#37)
}
